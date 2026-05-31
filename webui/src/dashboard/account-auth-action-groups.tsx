import type { ReactNode } from 'react';
import { KeyRound, LogIn, RefreshCw, Search, ShieldCheck, Smartphone } from 'lucide-react';
import { AccountActionRow, AccountActionRows, accountActionButton } from '@byte-v-forge/common-ui';
import type { ActionButtonDescriptor, AccountButtonActionSpec } from '@byte-v-forge/common-ui';
import { GPT_ACTIONS, type GptActionCatalog } from './action-catalog';
import { canLoginSession, canProbeAccount, canRegister, canUpdateWebAccessToken, loginActionHint, loginActionLabel, probeAccountHint } from './account-utils';
import type { Account } from './types';

export function AccountPrimaryActions(props: {
  account: Account;
  actionCatalog?: GptActionCatalog;
  busy: boolean;
  updatingWebAccessToken: boolean;
  onProbeAccount: (account: Account) => void;
  onRegister: (account: Account) => void;
  onRegisterProtocol: (account: Account) => void;
  onLogin: (account: Account) => void;
  onLoginProtocol: (account: Account) => void;
  onCodexOAuthAddPhone: (account: Account) => void;
  onCodexOAuthProtocol: (account: Account) => void;
  onUpdateWebAccessToken: (account: Account) => Promise<void>;
}) {
  const ctx = { catalog: props.actionCatalog, account: props.account, busy: props.busy, placement: 'account_header_browser' };
  return (
    <AccountActionRows className="accountPrimaryActionRows">
      <AccountActionRow className="accountPrimaryActionRow" buttonGroupClassName="sectionActions accountPrimaryActions" label="浏览器" actions={browserActions(props).map((spec) => accountActionButton(ctx, spec))} />
      <AccountActionRow className="accountPrimaryActionRow" buttonGroupClassName="sectionActions accountPrimaryActions" label="协议" actions={protocolActions(props)} />
      <AccountActionRow className="accountPrimaryActionRow" buttonGroupClassName="sectionActions accountPrimaryActions" label="工具" actions={toolActions(props)} />
    </AccountActionRows>
  );
}

function browserActions(props: Parameters<typeof AccountPrimaryActions>[0]): AccountButtonActionSpec<Account>[] {
  return [
    action('register', GPT_ACTIONS.register, '注册', <KeyRound size={14} />, canRegister, '账号不可注册', '浏览器注册', props.onRegister),
    action('login', GPT_ACTIONS.loginSession, '登录更新认证', <LogIn size={14} />, canLoginSession, '缺少邮箱', loginActionHint, props.onLogin),
    action('probe', GPT_ACTIONS.probeAccount, '探测账号', <Search size={14} />, canProbeAccount, '需要已注册账号', probeAccountHint, props.onProbeAccount),
  ];
}

function protocolActions(props: Parameters<typeof AccountPrimaryActions>[0]): ActionButtonDescriptor[] {
  const ctx = { catalog: props.actionCatalog, account: props.account, busy: props.busy, placement: 'account_header_protocol' };
  return [
    action('register-protocol', GPT_ACTIONS.registerProtocol, '协议注册', <ShieldCheck size={14} />, canRegister, '账号不可注册', '协议注册账号', props.onRegisterProtocol),
    action('login-protocol', GPT_ACTIONS.loginSessionProtocol, '协议登录', <LogIn size={14} />, canLoginSession, '缺少邮箱', loginActionLabel, props.onLoginProtocol),
    action('codex-oauth-protocol', GPT_ACTIONS.codexOAuthProtocol, 'Codex OAuth', <Smartphone size={14} />, canLoginSession, '缺少邮箱', '协议绑定 Codex OAuth', props.onCodexOAuthProtocol),
  ].map((spec) => accountActionButton(ctx, spec));
}

function toolActions(props: Parameters<typeof AccountPrimaryActions>[0]): ActionButtonDescriptor[] {
  const ctx = { catalog: props.actionCatalog, account: props.account, busy: props.busy, placement: 'account_header_tools' };
  return [
    {
      id: 'refresh-access-token',
      visible: canUpdateWebAccessToken(props.account),
      label: props.updatingWebAccessToken ? '更新中' : '更新 Web AT',
      hint: '使用当前 Session Token 获取 ChatGPT Web AT',
      icon: <RefreshCw size={14} />,
      disabled: props.busy || props.updatingWebAccessToken,
      onClick: () => void props.onUpdateWebAccessToken(props.account),
    },
    accountActionButton(ctx, action('codex-oauth-add-phone', GPT_ACTIONS.codexOAuth, 'Codex OAuth', <Smartphone size={14} />, canLoginSession, '缺少邮箱', '绑定 Codex OAuth / 手机', props.onCodexOAuthAddPhone)),
  ];
}

function action(id: string, actionID: string, fallbackLabel: string, icon: ReactNode, allowed: (account: Account) => boolean, disabledReason: string, hint: string | ((account: Account) => string), onClick: (account: Account) => void): AccountButtonActionSpec<Account> {
  return { id, actionID, fallbackLabel, icon, allowed, disabledReason, hint, onClick };
}
