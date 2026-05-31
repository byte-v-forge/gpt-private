import { Badge, EmptyBlock, KVList, api, formatUnix, useQuery } from '@byte-v-forge/common-ui';
import type { AccountProxyUsage } from './types';
import type { ListAccountProxyUsagesResponse } from '../proto/orchestrator_account';

export function AccountProxyHistoryPanel({ accountID }: { accountID: string }) {
  const query = useQuery({ queryKey: ['gpt', 'account-proxy-usages', accountID], queryFn: () => api<ListAccountProxyUsagesResponse>(`/api/gpt/accounts/${encodeURIComponent(accountID)}/proxy-usages?limit=100`), enabled: !!accountID });
  const rows = query.data?.usages || [];
  if (query.isLoading) return <EmptyBlock text="正在加载 IP 使用记录" />;
  if (query.error) return <EmptyBlock text="IP 使用记录加载失败" />;
  if (!rows.length) return <EmptyBlock text="暂无 IP 使用记录" />;
  return <section className="accountProxyHistoryPanel">{rows.map((row) => <ProxyUsageCard key={row.id} row={row} />)}</section>;
}

function ProxyUsageCard({ row }: { row: AccountProxyUsage }) {
  return (
    <article className="proxyUsageCard">
      <header className="proxyUsageHeader"><strong>{formatUnix(row.created_at)}</strong><StatusBadge row={row} /></header>
      <KVList items={[
        { id: 'purpose', label: '用途', value: row.purpose || '-' },
        { id: 'exit-ip', label: '出口 IP', value: row.exit_ip || '-', copyValue: row.exit_ip, mono: true },
        { id: 'location', label: '位置', value: [row.country_code, row.region, row.city].filter(Boolean).join(' / ') || '-' },
        { id: 'attempt', label: '尝试', value: String(row.attempt_index || '-') },
        { id: 'error', label: '错误', value: row.error_message, visible: !!row.error_message }
      ]} />
    </article>
  );
}

function StatusBadge({ row }: { row: AccountProxyUsage }) {
  const cls = row.accepted ? 'good' : row.error_message ? 'bad' : 'neutral';
  return <Badge className={`badge ${cls}`} variant="outline">{row.accepted ? '通过' : row.error_message ? '失败' : '记录'}</Badge>;
}
