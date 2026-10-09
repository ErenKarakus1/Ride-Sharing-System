import type { Location } from "../types";
import { Field } from "./Field";

type CoordinateEditorProps = {
  title: string;
  value: Location;
  onChange: (value: Location) => void;
};

export function CoordinateEditor({ title, value, onChange }: CoordinateEditorProps) {
  return (
    <div className="coordinate-editor">
      <h3>{title}</h3>
      <Field label="Latitude">
        <input
          type="number"
          step="0.000001"
          value={value.latitude}
          onChange={(event) => onChange({ ...value, latitude: Number(event.target.value) })}
        />
      </Field>
      <Field label="Longitude">
        <input
          type="number"
          step="0.000001"
          value={value.longitude}
          onChange={(event) => onChange({ ...value, longitude: Number(event.target.value) })}
        />
      </Field>
    </div>
  );
}
