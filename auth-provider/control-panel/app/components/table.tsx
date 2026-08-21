import React from "react";

export type Column<T> = {
  key: keyof T;
  header?: string;
  // render receives the cell value and the normalized row (keys lowercased)
  render?: (value: unknown, row: Record<string, unknown>) => React.ReactNode;
  width?: string | number;
};

export type TableProps<T> = {
  data: T[];
  columns?: Column<T>[];
  rowKey?: keyof T | ((row: T) => string);
  className?: string;
  empty?: React.ReactNode;
};

function keyToHeader(k: string) {
  // split camelCase / snake_case and capitalize words
  const words = k
    .replace(/([a-z0-9])([A-Z])/g, "$1 $2")
    .replace(/[_-]+/g, " ")
    .split(" ")
    .filter(Boolean)
    .map((w) => w[0].toUpperCase() + w.slice(1));
  return words.join(" ");
}

export function Table<T extends Record<string, unknown>>({
  data,
  columns,
  rowKey,
  className,
  empty = <div>No data</div>,
}: TableProps<T>) {
  // Normalize incoming data: create lowercase variants of keys so columns like 'name' match 'Name' or 'NAME'
  const { normalizedData, inferredColumns } = React.useMemo(() => {
    const nd: Record<string, unknown>[] = [];
    if (!data || data.length === 0)
      return { normalizedData: nd, inferredColumns: [] as Column<T>[] };

    for (const row of data as unknown as Record<string, unknown>[]) {
      const r: Record<string, unknown> = {};
      for (const k of Object.keys(row)) {
        const v = (row as Record<string, unknown>)[k];
        r[k] = v;
        const lk = k.toLowerCase();
        if (!(lk in r)) r[lk] = v;
      }
      nd.push(r);
    }

    if (columns && columns.length)
      return { normalizedData: nd, inferredColumns: columns as Column<T>[] };

    const firstKeys = Object.keys(nd[0]);
    const seen = new Set<string>();
    const cols: Column<T>[] = firstKeys
      .map((k) => k.toLowerCase())
      .filter((k) => {
        if (seen.has(k)) return false;
        seen.add(k);
        return true;
      })
      .map(
        (k) =>
          ({
            key: k as unknown as keyof T,
            header: keyToHeader(String(k)),
          }) as Column<T>,
      );

    return { normalizedData: nd, inferredColumns: cols };
  }, [columns, data]);

  if (!normalizedData || normalizedData.length === 0)
    return <div className={className}>{empty}</div>;

  return (
    <div className={className}>
      <table className="w-full table-auto border-collapse">
        <thead>
          <tr>
            {inferredColumns.map((col) => (
              <th
                key={String(col.key)}
                className="text-left p-2 border-b bg-gray-50"
                style={{ width: col.width }}
              >
                {col.header ?? keyToHeader(String(col.key))}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {normalizedData.map((nrow, i) => {
            const rk =
              typeof rowKey === "function"
                ? rowKey(nrow as unknown as T)
                : rowKey
                  ? String(nrow[String(rowKey)] ?? i)
                  : String(i);
            return (
              <tr key={String(rk)} className="border-b">
                {inferredColumns.map((col) => {
                  const keyStr = String(col.key);
                  const val = nrow[keyStr] ?? nrow[keyStr.toLowerCase()];
                  return (
                    <td key={String(col.key)} className="p-2 align-top">
                      {col.render ? col.render(val, nrow) : String(val ?? "")}
                    </td>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default Table;
