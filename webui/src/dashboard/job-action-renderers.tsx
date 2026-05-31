import type { ReactNode } from 'react';
import type { Job, WorkflowProgress } from './types';

export type WorkflowJobActionMessageKind = 'ok' | 'error';

export type WorkflowJobActionRendererProps = {
  job: Job;
  progress: WorkflowProgress | null;
  nowUnix: number;
  onChanged?: () => void | Promise<void>;
  onMessage?: (kind: WorkflowJobActionMessageKind, message: string) => void;
  onError?: (error: unknown) => void;
};

export type WorkflowJobActionRenderer = {
  id: string;
  statuses?: string[];
  render: (props: WorkflowJobActionRendererProps) => ReactNode;
};

const jobActionRenderers: WorkflowJobActionRenderer[] = [];

export function registerWorkflowJobActionRenderers(renderers: WorkflowJobActionRenderer[]) {
  for (const renderer of renderers) {
    const index = jobActionRenderers.findIndex((item) => item.id === renderer.id);
    if (index >= 0) jobActionRenderers[index] = renderer;
    else jobActionRenderers.push(renderer);
  }
}

export function renderWorkflowJobActions(props: WorkflowJobActionRendererProps) {
  return jobActionRenderers
    .filter((renderer) => !renderer.statuses?.length || renderer.statuses.includes(props.job.status))
    .map((renderer) => <div key={renderer.id}>{renderer.render(props)}</div>);
}
