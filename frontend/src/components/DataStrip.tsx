type DataStripProps = {
  items: Array<[string, string]>;
};

export function DataStrip({ items }: DataStripProps) {
  return (
    <div className="data-strip">
      {items.map(([label, value]) => (
        <div key={label}>
          <span>{label}</span>
          <strong>{value}</strong>
        </div>
      ))}
    </div>
  );
}
