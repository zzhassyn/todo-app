const TAG_PATTERN = /#([a-zA-Zа-яА-ЯёЁ0-9_-]{1,50})/g;
const PRIORITY_PATTERN = /!(low|medium|high)\b/i;

/**
 * Parses inline markup out of a task title string:
 *   "купить молоко #work !high" ->
 *     { title: "купить молоко", tags: ["work"], priority: "high" }
 *
 * Tags and the priority marker are stripped from the returned title.
 * Unrecognized "!word" sequences (not low/medium/high) are left as-is,
 * since they're probably just punctuation, not a priority marker.
 */
export function parseTaskInput(raw) {
  let priority = null;

  const priorityMatch = raw.match(PRIORITY_PATTERN);
  if (priorityMatch) {
    priority = priorityMatch[1].toLowerCase();
  }

  const tags = [];
  for (const match of raw.matchAll(TAG_PATTERN)) {
    const name = match[1].toLowerCase();
    if (!tags.includes(name)) tags.push(name);
  }

  const title = raw
    .replace(TAG_PATTERN, "")
    .replace(PRIORITY_PATTERN, "")
    .replace(/\s{2,}/g, " ")
    .trim();

  return { title, tags, priority };
}
