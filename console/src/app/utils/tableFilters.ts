export function searchItems<T>(items: T[], query: string, fields: (keyof T)[]): T[] {
  if (!query || query.trim() === '') {
    return items;
  }

  const lowerQuery = query.toLowerCase();
  return items.filter(item =>
    fields.some(field => {
      const value = item[field];
      return value != null && String(value).toLowerCase().includes(lowerQuery);
    })
  );
}

export function filterByStatus<T extends { status?: string | null }>(items: T[], status?: string | null): T[] {
  if (!status) {
    return items;
  }
  return items.filter(item => item.status === status);
}

export function filterByActionType<T extends { actionType?: string | null }>(
  items: T[],
  actionTypes?: string[] | null
): T[] {
  if (!actionTypes || actionTypes.length === 0) {
    return items;
  }
  return items.filter(item => item.actionType && actionTypes.includes(item.actionType));
}

export function filterByProvider<T extends { provider?: string | null }>(items: T[], providers?: string[] | null): T[] {
  if (!providers || providers.length === 0) {
    return items;
  }
  return items.filter(item => item.provider && providers.includes(item.provider));
}

export function paginateItems<T>(items: T[], page: number, perPage: number): T[] {
  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  return items.slice(startIndex, endIndex);
}

export function sortItems<T>(items: T[], field: keyof T, direction: 'asc' | 'desc'): T[] {
  return [...items].sort((a, b) => {
    const aVal = a[field];
    const bVal = b[field];

    if (typeof aVal === 'number' && typeof bVal === 'number') {
      return direction === 'asc' ? aVal - bVal : bVal - aVal;
    }

    const aStr = String(aVal || '');
    const bStr = String(bVal || '');
    return direction === 'asc' ? aStr.localeCompare(bStr) : bStr.localeCompare(aStr);
  });
}
