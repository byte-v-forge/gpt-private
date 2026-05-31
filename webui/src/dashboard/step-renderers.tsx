import type { ReactNode } from 'react';
import type { JobStep } from './types';

export type WorkflowStepRendererProps = {
  step: JobStep;
};

export type WorkflowStepRenderer = {
  id: string;
  stepNames: string[];
  label: string;
  render?: (props: WorkflowStepRendererProps) => ReactNode;
};

const stepRenderers: WorkflowStepRenderer[] = [];

export function registerWorkflowStepRenderers(renderers: WorkflowStepRenderer[]) {
  for (const renderer of renderers) {
    const index = stepRenderers.findIndex((item) => item.id === renderer.id);
    if (index >= 0) stepRenderers[index] = renderer;
    else stepRenderers.push(renderer);
  }
}

export function workflowStepRenderer(stepName: string) {
  return stepRenderers.find((renderer) => renderer.stepNames.includes(stepName));
}

export function renderWorkflowStep(step: JobStep) {
  const renderer = workflowStepRenderer(step.step_name);
  return renderer?.render ? renderer.render({ step }) : null;
}
