import type { ActionKey, ActionState } from "../types";

export const actionKeys: ActionKey[] = [
  "auth",
  "fare",
  "ride",
  "payment",
  "driver-location",
  "matching",
  "lifecycle",
];

export function initialActionStates(): Record<ActionKey, ActionState> {
  return Object.fromEntries(
    actionKeys.map((key) => [key, { loading: false, message: "", error: "" }]),
  ) as Record<ActionKey, ActionState>;
}
