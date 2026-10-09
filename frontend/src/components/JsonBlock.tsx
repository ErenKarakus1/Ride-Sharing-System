type JsonBlockProps = {
  title: string;
  value: unknown;
};

export function JsonBlock({ title, value }: JsonBlockProps) {
  return (
    <div className="json-block">
      <h3>{title}</h3>
      <pre>{value ? JSON.stringify(value, null, 2) : "null"}</pre>
    </div>
  );
}
