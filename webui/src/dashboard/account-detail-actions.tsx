import { AccountActionRow, AccountActionRows, AccountOTPActionRow, accountActionAvailability, accountDeleteButtonAction, accountMailboxInboxHintForAccount, hasVisibleAction } from '@byte-v-forge/common-ui';
import type { ActionButtonDescriptor } from '@byte-v-forge/common-ui';
import { PaymentChannelIcon } from './account-badges';
import { GPT_ACTIONS, type GptActionCatalog } from './action-catalog';
import { canGoPayPayment } from './account-utils';
import { GO_PAY_PAYMENT_CHANNELS, goPayPaymentActionLabel } from './gopay-utils';
import type { AccountMailboxContext } from '@byte-v-forge/common-ui';
import type { Account, ConcreteGoPayPaymentChannel, LatestOtp } from './types';

export function AccountDetailActions({ account, actionCatalog, showSecrets, busy, inboxLoading, mailboxContext, latestOtp, canFetchOTP, goPayAvailable, goPayUnavailableReason, onCopy, onFetchInbox, onGoPayPayment }: {
  account: Account;
  actionCatalog?: GptActionCatalog;
  showSecrets: boolean;
  busy: boolean;
  inboxLoading: boolean;
  mailboxContext: AccountMailboxContext | null;
  latestOtp: LatestOtp | null;
  canFetchOTP: boolean;
  goPayAvailable: boolean;
  goPayUnavailableReason: string;
  onCopy: (label: string, value: string) => void;
  onFetchInbox: (account: Account) => Promise<void>;
  onGoPayPayment: (account: Account, channel: ConcreteGoPayPaymentChannel) => void;
}) {
  const channelRow = channelActions(actionCatalog, account, busy, goPayAvailable, goPayUnavailableReason, onGoPayPayment);
  return (
    <AccountActionRows>
      <AccountOTPActionRow
        latestOtp={latestOtp}
        showSecrets={showSecrets}
        canRefresh={canFetchOTP}
        refreshDisabled={busy || inboxLoading}
        refreshHint={accountMailboxInboxHintForAccount(account, mailboxContext, showSecrets)}
        onCopy={onCopy}
        onRefresh={() => void onFetchInbox(account)}
      />
      {hasVisibleAction(channelRow) && <AccountActionRow label="激活" actions={channelRow} />}
    </AccountActionRows>
  );
}

export function AccountDangerActions({ account, busy, onDelete }: {
  account: Account;
  busy: boolean;
  onDelete: (account: Account) => Promise<void>;
}) {
  return (
    <AccountActionRows className="bottomActionRows">
      <AccountActionRow label="危险" actions={dangerActions(account, busy, onDelete)} />
    </AccountActionRows>
  );
}

function channelActions(catalog: GptActionCatalog | undefined, account: Account, busy: boolean, goPayAvailable: boolean, goPayUnavailableReason: string, onGoPayPayment: (account: Account, channel: ConcreteGoPayPaymentChannel) => void): ActionButtonDescriptor[] {
  return GO_PAY_PAYMENT_CHANNELS.map((channel) => {
    const availability = accountActionAvailability(catalog, GPT_ACTIONS.goPayPayment, account, 'account_detail');
    return {
      id: `gopay-payment-${channel}`,
      visible: availability.visible,
      label: goPayPaymentActionLabel(channel),
      hint: !goPayAvailable ? goPayUnavailableReason : availability.reason || (canGoPayPayment(account) ? 'GoPay 激活支付渠道' : '需要已注册且未激活账号'),
      icon: <span className="activationPaymentIcon"><PaymentChannelIcon channel={channel} /></span>,
      className: 'activationAction',
      disabled: busy || !goPayAvailable || !availability.enabled || !canGoPayPayment(account),
      onClick: () => onGoPayPayment(account, channel),
    };
  });
}

function dangerActions(account: Account, busy: boolean, onDelete: (account: Account) => Promise<void>): ActionButtonDescriptor[] {
  return [accountDeleteButtonAction(() => onDelete(account), busy)];
}
