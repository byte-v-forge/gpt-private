import { Zap } from 'lucide-react';
import {
  AccountCarrierList,
  RecordActionButtons,
  RecordActions,
  accountDeleteRowAction,
  accountActionAvailability,
  accountCarrierID,
  type AccountListPagination
} from '@byte-v-forge/common-ui';
import type { RowActionDescriptor } from '@byte-v-forge/common-ui';
import {
  canGoPayPayment,
  isInvalidGptAccount,
  isUserAlreadyExistsAccount
} from './account-utils';
import { accountActivationChannel, accountCodexPhoneState } from './account-job-semantics';
import { AccountChannelTag, AccountCodexPhoneTag, AccountSignalBadge, PaymentChannelIcon } from './account-badges';
import { GPT_ACTIONS, type GptActionCatalog } from './action-catalog';
import { gptAccountRecord } from './account-record';
import { AccountRowAuthGroups } from './account-row-auth-groups';
import { OpenAIIcon } from './brand-icons';
import { GO_PAY_PAYMENT_CHANNELS, goPayPaymentActionLabel } from './gopay-utils';
import type { Account, ConcreteGoPayPaymentChannel, Job } from './types';

export function AccountTable({ accounts, jobs, selected, actionCatalog, showSecrets, runningAccountIds, runningWorkflowByAccountID, pagination, busy, goPayAvailable, goPayUnavailableReason, onSelect, onRegisterProtocol, onGoPayPayment, onDelete }: {
  accounts: Account[];
  jobs: Job[];
  selected?: string;
  actionCatalog?: GptActionCatalog;
  showSecrets: boolean;
  runningAccountIds: Set<string>;
  runningWorkflowByAccountID: Map<string, Job>;
  pagination?: AccountListPagination;
  busy: boolean;
  goPayAvailable: boolean;
  goPayUnavailableReason: string;
  onSelect: (a: Account) => void;
  onRegisterProtocol: (a: Account) => void;
  onGoPayPayment: (a: Account, channel: ConcreteGoPayPaymentChannel) => void;
  onDelete: (a: Account) => void | Promise<void>;
}) {
  const byID = new Map(accounts.map((account) => [accountCarrierID(account), account] as const));
  return (
    <AccountCarrierList
      carriers={accounts}
      selectedID={selected}
      emptyText="暂无账号。可以先创建账号，或切换为全部状态查看。"
      listClassName="accountsList"
      pagination={pagination}
      onSelectCarrier={onSelect}
      recordOf={(account) => gptAccountRecord(account, showSecrets)}
      config={{
        icon: () => <OpenAIIcon size={15} />,
        title: (record) => <span className="accountCardEmail" title={record.display_name}>{record.display_name}</span>,
        subtitle: () => '',
        meta: (record) => {
          const account = byID.get(accountCarrierID(record));
          if (!account) return null;
          const activationChannel = accountActivationChannel(account, jobs, actionCatalog);
          const phoneState = accountCodexPhoneState(account, jobs, actionCatalog);
          return (
            <div className="accountCardTags">
              <AccountSignalBadge account={account} compact />
              <AccountCodexPhoneTag state={phoneState} />
              <AccountChannelTag channel={activationChannel} />
            </div>
          );
        }
      }}
      renderChildren={(account) => {
        const accountID = accountCarrierID(account);
        const accountBusy = runningAccountIds.has(accountID);
        const currentWorkflow = runningWorkflowByAccountID.get(accountID);
        return (
          <AccountRowActions
            account={account}
            actionCatalog={actionCatalog}
            accountBusy={accountBusy}
            currentWorkflow={currentWorkflow}
            busy={busy}
            goPayAvailable={goPayAvailable}
            goPayUnavailableReason={goPayUnavailableReason}
            onRegisterProtocol={onRegisterProtocol}
            onGoPayPayment={onGoPayPayment}
            onDelete={onDelete}
          />
        );
      }}
    />
  );
}

function AccountRowActions({ account, actionCatalog, accountBusy, currentWorkflow, busy, goPayAvailable, goPayUnavailableReason, onRegisterProtocol, onGoPayPayment, onDelete }: {
  account: Account;
  actionCatalog?: GptActionCatalog;
  accountBusy: boolean;
  currentWorkflow?: Job;
  busy: boolean;
  goPayAvailable: boolean;
  goPayUnavailableReason: string;
  onRegisterProtocol: (a: Account) => void;
  onGoPayPayment: (a: Account, channel: ConcreteGoPayPaymentChannel) => void;
  onDelete: (a: Account) => void | Promise<void>;
}) {
  if (isInvalidGptAccount(account)) {
    const actions: RowActionDescriptor[] = [accountDeleteRowAction(() => onDelete(account), busy)];
    return (
      <RecordActions className="rowActions">
        <div className="rowActionsMain"><RecordActionButtons actions={actions} /></div>
      </RecordActions>
    );
  }

  if (accountBusy && currentWorkflow && !isUserAlreadyExistsAccount(account)) {
    return (
      <RecordActions className="rowActions">
        <div className="rowActionsMain"><span className="accountWorkflowNotice">流程运行中，请到工作流页处理</span></div>
      </RecordActions>
    );
  }

  const payment = accountActionAvailability(actionCatalog, GPT_ACTIONS.goPayPayment, account, 'account_row');
  const paymentActions: RowActionDescriptor[] = GO_PAY_PAYMENT_CHANNELS.filter((channel) => channel !== 'wa' && payment.visible).map((channel) => ({
    label: !goPayAvailable && goPayUnavailableReason ? `${goPayPaymentActionLabel(channel)} · ${goPayUnavailableReason}` : goPayPaymentActionLabel(channel),
    icon: <span className="activationPaymentIcon"><Zap size={13} /><PaymentChannelIcon channel={channel} /></span>,
    onClick: () => onGoPayPayment(account, channel),
    disabled: busy || !goPayAvailable || !payment.enabled || !canGoPayPayment(account),
    kind: 'secondary' as const,
    className: 'paymentIconAction activationAction'
  }));
  const leftActions = paymentActions;
  return (
    <RecordActions className="rowActions">
      <div className="rowActionsMain splitRowActions">
        <div className="rowActionsLeft"><RecordActionButtons actions={leftActions} /></div>
        <div className="rowActionsRight">
          <AccountRowAuthGroups account={account} actionCatalog={actionCatalog} busy={busy} onRegisterProtocol={onRegisterProtocol} />
        </div>
      </div>
    </RecordActions>
  );
}
