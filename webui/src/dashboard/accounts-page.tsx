import { useState } from 'react';
import { accountCarrierEmail, accountCarrierID, accountMailboxContextForEmail, DetailDrawer, latestOtpForEmail, ToastMessage, WorkspaceTabbedPanel } from '@byte-v-forge/common-ui';
import { AccountDetails } from './account-details';
import { GPT_ACTIONS, useGptActionCatalog, type GptActionCatalog } from './action-catalog';
import { useGptAccountActions } from './account-actions';
import { useGptAccountData } from './account-data';
import { useGptAccountEventCache } from './account-events';
import { accountActivationChannel, accountCodexPhoneState } from './account-job-semantics';
import { GptAccountsView } from './view';

type GptPageTab = 'accounts';

export function GptAccountsPage() {
  const actionCatalog = useGptActionCatalog();

  return (
    <WorkspaceTabbedPanel<GptPageTab>
      title="GPT"
      defaultValue="accounts"
      tabsListVariant="line"
      tabs={[
        {
          value: 'accounts',
          label: '账号',
          content: <GptAccountsTab actionCatalog={actionCatalog.data} />,
          contentClassName: 'flex flex-col overflow-hidden'
        }
      ]}
    />
  );
}

function GptAccountsTab({ actionCatalog }: { actionCatalog?: GptActionCatalog }) {
  const [selectedAccountID, setSelectedAccountID] = useState('');
  const [showSecrets, setShowSecrets] = useState(true);
  const data = useGptAccountData(selectedAccountID, setSelectedAccountID);
  const actions = useGptAccountActions(data, showSecrets, setSelectedAccountID, actionCatalog);
  const busy = actions.working || data.busy;
  const selectedPhoneState = data.selected ? accountCodexPhoneState(data.selected, data.jobs, actions.actionCatalog) : null;
  useGptAccountEventCache();

  return (
    <>
      <ToastMessage toast={actions.toast.toast} />
      <GptAccountsView
        accounts={data.accounts}
        jobs={data.jobs}
        selectedAccountId={selectedAccountID}
        actionCatalog={actions.actionCatalog}
        showSecrets={showSecrets}
        busy={busy}
        cleaningInvalidAccounts={actions.cleaningInvalidAccounts}
        runningAccountIds={data.runningIds}
        runningWorkflowByAccountID={data.runningByAccount}
        accountsPagination={data.accountsPagination}
        onCreateDone={async (message) => {
          actions.toast.showOK(message);
          await data.invalidate();
        }}
        onError={actions.toast.showError}
        onToggleSecrets={() => setShowSecrets((value) => !value)}
        onCleanInvalidAccounts={actions.cleanInvalidAccounts}
        onSelectAccount={(account) => setSelectedAccountID(accountCarrierID(account))}
        onRegisterProtocol={(account) => actions.runWorkflow(GPT_ACTIONS.registerProtocol, account)}
        onCodexOAuthBatchAddPhone={actions.runCodexOAuthBatchAddPhone}
        goPayAvailable={actions.goPayAvailable}
        goPayUnavailableReason={actions.goPayAvailabilityReason}
        onGoPayPayment={(account, channel) => void actions.runGoPayPayment(account, channel)}
        onDeleteAccount={actions.deleteAccount}
      />
      <DetailDrawer open={!!data.selected} title="GPT账号详情" size="wide" onClose={() => setSelectedAccountID('')}>
        {data.selected && (
          <AccountDetails
            account={data.selected}
            jobs={data.jobs}
            actionCatalog={actions.actionCatalog}
            showSecrets={showSecrets}
            busy={busy}
            inboxLoading={actions.inboxLoading}
            mailboxContext={accountMailboxContextForEmail(data.mailboxes, data.allocations, data.selected)}
            latestOtp={latestOtpForEmail(actions.inbox, data.mailboxes, accountCarrierEmail(data.selected))}
            activationChannel={accountActivationChannel(data.selected, data.jobs, actions.actionCatalog)}
            codexPhoneState={selectedPhoneState!}
            updatingWebAccessToken={actions.updatingWebAccessTokens.has(accountCarrierID(data.selected))}
            onCopy={actions.toast.copyValue}
            onFetchInbox={actions.fetchInbox}
            onSessionSave={(account, sessionToken) => actions.updateAccount(account, { session_token: sessionToken }, '认证信息已更新')}
            onAccessSave={(account, accessToken) => actions.updateAccount(account, { access_token: accessToken }, '认证信息已更新')}
            onProbeAccount={(account) => actions.runWorkflow(GPT_ACTIONS.probeAccount, account)}
            onRegister={(account) => actions.runWorkflow(GPT_ACTIONS.register, account)}
            onRegisterProtocol={(account) => actions.runWorkflow(GPT_ACTIONS.registerProtocol, account)}
            onLogin={(account) => actions.runWorkflow(GPT_ACTIONS.loginSession, account)}
            onLoginProtocol={(account) => actions.runWorkflow(GPT_ACTIONS.loginSessionProtocol, account)}
            onCodexOAuthAddPhone={(account) => actions.runWorkflow(GPT_ACTIONS.codexOAuth, account)}
            onCodexOAuthProtocol={(account) => actions.runWorkflow(GPT_ACTIONS.codexOAuthProtocol, account)}
            goPayAvailable={actions.goPayAvailable}
            goPayUnavailableReason={actions.goPayAvailabilityReason}
            onGoPayPayment={(account, channel) => void actions.runGoPayPayment(account, channel)}
            onUpdateWebAccessToken={actions.updateWebAccessToken}
            onDelete={actions.deleteAccount}
          />
        )}
      </DetailDrawer>
    </>
  );
}
