import { parseISO, format, formatDistanceToNow } from 'date-fns';
import { createParser } from 'nuqs';

export function parseScanTimestamp(ts: string | undefined) {
  if (!ts || ts.startsWith('0001-01-01')) return 'Never';
  const date = parseISO(ts);
  const absolute = format(date, 'HH:mm:ss - dd/MM/yyyy');
  const relative = formatDistanceToNow(date, { addSuffix: true });
  return `${absolute} (${relative})`;
}

export function parseNumberToCurrency(value: number | undefined) {
  if (value === undefined || value === null) return '$0.00';
  return value.toLocaleString('en-US', {
    style: 'currency',
    currency: 'USD',
  });
}

export function resolveResourcePath(resourceType: string, resourceName: string): string {
  if (resourceType === 'Cluster') {
    return `/clusters/${resourceName}`;
  }

  if (resourceType === 'Instance') {
    return `/instances/${resourceName}`;
  }

  if (resourceType === 'Account') {
    return `/accounts/${resourceName}`;
  }

  return '#';
}

// Nullable boolean: "true" -> true, "false" -> false, missing/other -> null
export const parseAsBooleanNullable = createParser<boolean | null>({
  parse: value => {
    if (value === 'true') return true;
    if (value === 'false') return false;
    return null;
  },
  serialize: value => {
    // nuqs expects a string; return empty string to represent "unset"
    if (value === null) return '';
    return value ? 'true' : 'false';
  },
});
