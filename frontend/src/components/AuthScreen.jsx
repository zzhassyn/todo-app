import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import { ApiError } from "../api/client";

const initialRegister = { fullName: "", email: "", password: "" };
const initialLogin = { email: "", password: "" };

function describeError(err, mode) {
  if (!(err instanceof ApiError)) {
    return "не получилось подключиться к серверу — проверьте, что бэкенд запущен";
  }
  if (mode === "register" && err.status === 409) {
    return "этот email уже зарегистрирован";
  }
  if (mode === "login" && err.status === 401) {
    return "неверный email или пароль";
  }
  if (err.status === 400) {
    return "проверьте корректность введённых данных";
  }
  return err.message;
}

export default function AuthScreen() {
  const { login, register } = useAuth();
  const [mode, setMode] = useState("login"); // login | register
  const [form, setForm] = useState(initialLogin);
  const [error, setError] = useState(null);
  const [busy, setBusy] = useState(false);

  function switchMode(next) {
    setMode(next);
    setForm(next === "login" ? initialLogin : initialRegister);
    setError(null);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError(null);
    setBusy(true);
    try {
      if (mode === "login") {
        await login({ email: form.email, password: form.password });
      } else {
        await register({
          full_name: form.fullName,
          email: form.email,
          password: form.password,
        });
      }
    } catch (err) {
      setError(describeError(err, mode));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="auth-screen">
      <div className="auth-card">
        <div className="auth-card__head">
          <span className="auth-card__dot" aria-hidden="true" />
          <span className="auth-card__path">
            ~/todo/{mode === "login" ? "session.open" : "session.new"}
          </span>
        </div>

        <div className="auth-card__tabs" role="tablist" aria-label="режим входа">
          <button
            type="button"
            role="tab"
            aria-selected={mode === "login"}
            className={`auth-tab ${mode === "login" ? "is-active" : ""}`}
            onClick={() => switchMode("login")}
          >
            войти
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={mode === "register"}
            className={`auth-tab ${mode === "register" ? "is-active" : ""}`}
            onClick={() => switchMode("register")}
          >
            создать аккаунт
          </button>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          {mode === "register" && (
            <label className="field">
              <span className="field__label">имя</span>
              <input
                type="text"
                required
                minLength={3}
                maxLength={100}
                value={form.fullName}
                autoComplete="name"
                onChange={(e) => setForm((f) => ({ ...f, fullName: e.target.value }))}
                placeholder="Ада Лавлейс"
              />
            </label>
          )}

          <label className="field">
            <span className="field__label">email</span>
            <input
              type="email"
              required
              value={form.email}
              autoComplete="email"
              onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
              placeholder="ada@example.com"
            />
          </label>

          <label className="field">
            <span className="field__label">пароль</span>
            <input
              type="password"
              required
              minLength={8}
              maxLength={72}
              value={form.password}
              autoComplete={mode === "login" ? "current-password" : "new-password"}
              onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
              placeholder="минимум 8 символов"
            />
          </label>

          {error && (
            <p className="auth-error">
              <span aria-hidden="true">! </span>
              {error}
            </p>
          )}

          <button type="submit" className="btn btn--primary" disabled={busy}>
            {busy ? "выполняется…" : mode === "login" ? "войти" : "создать аккаунт"}
          </button>
        </form>

        <p className="auth-card__footnote">
          {mode === "login" ? (
            <>
              ещё нет аккаунта?{" "}
              <button type="button" className="link-btn" onClick={() => switchMode("register")}>
                создать
              </button>
            </>
          ) : (
            <>
              уже есть аккаунт?{" "}
              <button type="button" className="link-btn" onClick={() => switchMode("login")}>
                войти
              </button>
            </>
          )}
        </p>
      </div>
    </div>
  );
}
