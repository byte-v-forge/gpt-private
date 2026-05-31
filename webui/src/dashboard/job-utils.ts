import { formatUnix, objectValue } from '@byte-v-forge/common-ui';
import type { JobStep } from '../proto/orchestrator_job';

export function formatJobTime(value: number) {
  return formatUnix(value);
}

export function stepDetailData(step?: Pick<JobStep, 'detail'> | null) {
  const detail = objectValue(step?.detail);
  for (const value of Object.values(detail) as unknown[]) {
    if (value && typeof value === 'object') return value as Record<string, any>;
  }
  return detail;
}
