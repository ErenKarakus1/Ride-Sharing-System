import type { Location, LocationPreset } from "../types";
import { Field } from "./Field";

type CoordinateEditorProps = {
  title: string;
  value: Location;
  presets?: LocationPreset[];
  onChange: (value: Location) => void;
};

export function CoordinateEditor({ title, value, presets = [], onChange }: CoordinateEditorProps) {
  return (
    <div className="coordinate-editor">
      <h3>{title}</h3>
      {presets.length > 0 && (
        <div className="preset-row" aria-label={`${title} presets`}>
          {presets.map((preset) => (
            <button
              className="preset"
              key={preset.label}
              onClick={() => onChange(preset.location)}
              type="button"
            >
              {preset.label}
            </button>
          ))}
        </div>
      )}
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
