import type { Account, GPTEmailAllocation } from '../proto/gpt_account';
import type { FetchAccountMailboxResponse as ProtoFetchAccountMailboxResponse } from '../proto/orchestrator_account';
import type { InboxMessage, InboxResponse, InboxResult, LatestOtp, Mailbox, MailboxDomain, MailboxProviderCapability } from '@byte-v-forge/common-ui';
import type { Job, JobData, JobSnapshot, JobStep, WorkflowProgress } from '../proto/orchestrator_job';
import type { ProxyEdgeAccessCheck, ProxyIPFraudCheck } from '@byte-v-forge/common-ui/proto/byte/v/forge/contracts/proxyruntime/v1/proxy_runtime';

export type FetchAccountMailboxResponse = Omit<ProtoFetchAccountMailboxResponse, 'inbox'> & { inbox?: InboxResponse };

export type { Account, GPTEmailAllocation, InboxMessage, InboxResponse, InboxResult, Job, JobData, JobSnapshot, JobStep, LatestOtp, Mailbox, MailboxDomain, MailboxProviderCapability, WorkflowProgress };

export type GoPayOTPChannel = '' | 'sms' | 'wa';
export type GoPayPaymentChannel = GoPayOTPChannel | 'app_wa' | 'account';
export type ConcreteGoPayPaymentChannel = Exclude<GoPayPaymentChannel, ''>;
export type DisplayLabelMap = Record<string, string>;

export type AccountBrowserFingerprint = {
  account_id: string;
  country_code: string;
  region: string;
  browser_profile_template: string;
  browser_family: string;
  browser_major_version: string;
  os_family: string;
  tls_profile_family: string;
  tls_fingerprint_variant: string;
  locale: string;
  timezone: string;
  user_agent: string;
  accept_language: string;
  language: string;
  device_id: string;
  created_at?: number;
  updated_at?: number;
};


export type AccountProxyChainHop = {
  hop_id: string;
  order: number;
  role: string;
  source_kind: string;
  source_id: string;
  source_display_name: string;
  node_id: string;
  node_display_name: string;
  provider_id: string;
  gateway_id: string;
  gateway_display_name: string;
  observed_ip: string;
  country_code: string;
  region: string;
  city: string;
  status: string;
  delay_ms: number;
};

export type AccountProxyChain = {
  chain_id: string;
  hops: AccountProxyChainHop[];
};

export type AccountProxyUsage = {
  id: string;
  job_id: string;
  n8n_execution_id: string;
  purpose: string;
  proxy_url_hash: string;
  session_id_hash: string;
  exit_ip: string;
  country_code: string;
  region: string;
  city: string;
  ip_fraud_check?: Partial<ProxyIPFraudCheck>;
  edge_access_check?: Partial<ProxyEdgeAccessCheck>;
  target_reachable: boolean;
  attempt_index: number;
  accepted: boolean;
  error_message: string;
  created_at: number;
  chain?: AccountProxyChain;
};
