import { ExternalLink } from 'lucide-react';
import { Badge, Button, EmptyBlock, accountActionLabel, formatUnix, statusText } from '@byte-v-forge/common-ui';
import type { GptActionCatalog } from './action-catalog';
import type { Job } from './types';

export function AccountWorkflowTimelinePanel({ accountID, jobs, actionCatalog }: { accountID: string; jobs: Job[]; actionCatalog?: GptActionCatalog }) {
  const current = jobs.filter((job) => job.account_id === accountID).sort((a, b) => (b.updated_at || 0) - (a.updated_at || 0))[0];
  if (!current) return <EmptyBlock text="暂无工作流记录" />;
  return <section className="accountWorkflowTimelinePanel"><WorkflowRunLine job={current} actionCatalog={actionCatalog} /></section>;
}

function WorkflowRunLine({ job, actionCatalog }: { job: Job; actionCatalog?: GptActionCatalog }) {
  const workflow = accountActionLabel(actionCatalog, job.action, compact(job.action));
  return (
    <article className="workflowRunLine">
      <StatusBadge status={job.status} />
      <span className="workflowRunLineItem">当前流程：<strong title={workflow}>{workflow}</strong></span>
      <span className="workflowRunLineItem">当前步骤：<strong title={job.last_step}>{compact(job.last_step || '-')}</strong></span>
      <span className={`workflowRunMeta ${job.error_message ? 'workflowRunMetaError' : ''}`}>{job.error_message || `${formatUnix(job.updated_at)} · ${job.n8n_execution_id ? `n8n #${job.n8n_execution_id}` : job.job_id}`}</span>
      <Button asChild variant="outline" size="sm"><a href={workflowRuntimeHref(job)}><ExternalLink />工作流页</a></Button>
    </article>
  );
}

function workflowRuntimeHref(job: Job) {
  const params = new URLSearchParams();
  if (job.n8n_execution_id) params.set('execution_id', job.n8n_execution_id);
  if (job.job_id) params.set('run_id', job.job_id);
  return `/workflow-runtime/workflow?${params.toString()}`;
}

function StatusBadge({ status }: { status?: string }) {
  const normalized = (status || 'UNKNOWN').toUpperCase();
  return <Badge className={`badge ${statusClass(normalized)}`} variant="outline">{statusText(normalized)}</Badge>;
}

function statusClass(status?: string) {
  const normalized = (status || '').toUpperCase();
  if (normalized === 'SUCCEEDED') return 'good';
  if (normalized === 'RUNNING' || normalized === 'CREATED') return 'mid';
  if (normalized === 'CANCELED' || normalized.startsWith('FAILED')) return 'bad';
  return 'neutral';
}

function compact(value?: string) {
  return (value || '').replaceAll('_', ' ').toLowerCase().replace(/(^|\s)\S/g, (s) => s.toUpperCase());
}
