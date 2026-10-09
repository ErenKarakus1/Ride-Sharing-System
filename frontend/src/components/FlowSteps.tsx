import { CheckCircle2, Circle } from "lucide-react";

export type FlowStep = {
  label: string;
  done: boolean;
  active?: boolean;
};

type FlowStepsProps = {
  steps: FlowStep[];
};

export function FlowSteps({ steps }: FlowStepsProps) {
  return (
    <div className="flow-steps">
      {steps.map((step) => {
        const Icon = step.done ? CheckCircle2 : Circle;
        return (
          <div className={stepClass(step)} key={step.label}>
            <Icon size={16} />
            <span>{step.label}</span>
          </div>
        );
      })}
    </div>
  );
}

function stepClass(step: FlowStep) {
  if (step.done) return "flow-step done";
  if (step.active) return "flow-step active";
  return "flow-step";
}
