import { Eye, EyeOff, Phone, Trash2 } from 'lucide-react';
import {
  PanelHeader,
  ToolbarIconButton,
  accountActionAvailability,
  accountActionLabel,
  type AccountListPagination,
} from '@byte-v-forge/common-ui';
import type { Account, ConcreteGoPayPaymentChannel, Job } from './types';
import { AccountTable } from './account-table';
import { CreateAccountForm } from './create-account';
import { canLoginSession, isInvalidGptAccount } from './account-utils';
import { accountCodexPhoneState } from './account-job-semantics';
import { GPT_ACTIONS, type GptActionCatalog } from './action-catalog';
import { invalidAccountsForCleanup } from './account-cleanup-actions';
import { OpenAIIcon } from './brand-icons';

export type GptAccountsViewProps = {
  accounts: Account[];
  jobs: Job[];
  selectedAccountId?: string;
  actionCatalog?: GptActionCatalog;
  showSecrets: boolean;
  busy: boolean;
  cleaningInvalidAccounts: boolean;
  runningAccountIds: Set<string>;
  runningWorkflowByAccountID: Map<string, Job>;
  accountsPagination?: AccountListPagination;
  onCreateDone: (message: string) => Promise<void>;
  onError: (message: string) => void;
  onToggleSecrets: () => void;
  onCleanInvalidAccounts: () => void | Promise<void>;
  onSelectAccount: (account: Account) => void;
  onRegisterProtocol: (account: Account) => void | Promise<void>;
  onCodexOAuthBatchAddPhone: (accounts: Account[]) => void | Promise<void>;
  goPayAvailable: boolean;
  goPayUnavailableReason: string;
  onGoPayPayment: (account: Account, channel: ConcreteGoPayPaymentChannel) => void;
  onDeleteAccount: (account: Account) => void | Promise<void>;
};

export function GptAccountsView(props: GptAccountsViewProps) {
  const addPhoneAccounts = props.accounts.filter((account) => canLoginSession(account) && !accountCodexPhoneState(account, props.jobs, props.actionCatalog).confirmed);
  const batchAddPhone = accountActionAvailability(props.actionCatalog, GPT_ACTIONS.codexOAuthBatchAddPhone, undefined, 'account_bulk');
  const invalidAccounts = invalidAccountsForCleanup(props.accounts);
  return (
    <>
      <PanelHeader title="GPT账号" icon={<OpenAIIcon size={16} />}>
        <div className="headerControls accountHeaderControls">
          <CreateAccountForm compact onDone={props.onCreateDone} onError={props.onError} />
          {addPhoneAccounts.length > 0 && batchAddPhone.visible && (
            <ToolbarIconButton
              label={`${accountActionLabel(props.actionCatalog, GPT_ACTIONS.codexOAuthBatchAddPhone, '批量 Add Phone', 'account_bulk')} · ${addPhoneAccounts.length} 个未加手机账号`}
              icon={<Phone size={15} />}
              disabled={props.busy || !batchAddPhone.enabled}
              onClick={() => void props.onCodexOAuthBatchAddPhone(addPhoneAccounts)}
            />
          )}
          {invalidAccounts.length > 0 && (
            <ToolbarIconButton
              label={props.cleaningInvalidAccounts ? '清理中' : `清理失效账号 · ${invalidAccounts.length}`}
              icon={<Trash2 size={15} />}
              disabled={props.busy || props.cleaningInvalidAccounts}
              onClick={() => void props.onCleanInvalidAccounts()}
            />
          )}
          <ToolbarIconButton label={props.showSecrets ? '隐藏敏感信息' : '显示敏感信息'} icon={props.showSecrets ? <EyeOff size={15} /> : <Eye size={15} />} onClick={props.onToggleSecrets} />
        </div>
      </PanelHeader>
      <AccountTable
        accounts={props.accounts}
        jobs={props.jobs}
        selected={props.selectedAccountId}
        actionCatalog={props.actionCatalog}
        showSecrets={props.showSecrets}
        runningAccountIds={props.runningAccountIds}
        runningWorkflowByAccountID={props.runningWorkflowByAccountID}
        pagination={props.accountsPagination}
        busy={props.busy}
        onSelect={props.onSelectAccount}
        onRegisterProtocol={props.onRegisterProtocol}
        goPayAvailable={props.goPayAvailable}
        goPayUnavailableReason={props.goPayUnavailableReason}
        onGoPayPayment={props.onGoPayPayment}
        onDelete={props.onDeleteAccount}
      />
    </>
  );
}
