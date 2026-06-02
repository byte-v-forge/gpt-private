import { useEffect, useState } from 'react';
import { ExternalLink, WalletCards } from 'lucide-react';
import QRCode from 'qrcode';
import {
  AccountActionRow,
  AccountRowActionGroups,
  Button,
  accountActionAvailability,
  accountActionLabel,
  accountCarrierID,
  api,
  buttonHint,
  hasVisibleAction,
  objectValue,
  stringValue,
  useAsyncActionRunner,
  useQuery,
  type ActionButtonDescriptor,
} from '@byte-v-forge/common-ui';
import type { RowActionDescriptor } from '@byte-v-forge/common-ui';
import { accountIsActivated, hasRegisteredSession, isInvalidGptAccount, isUserAlreadyExistsAccount } from './account-utils';
import { jobDataObject } from './job-data';
import type { GptActionCatalog, GptActionPlacement } from './action-catalog';
import type { AccountWorkflowRunner } from './account-action-specs';
import type { WorkflowJobActionRendererProps } from './job-action-renderers';
import type { Account, Job, WorkflowProgress } from './types';

const PRIVATE_OWNER = 'gpt-private';
const PAYMENT_CAPABILITY = 'payment';

type GoPayAvailability = { available: boolean; loading: boolean; reason: string };

export function GoPayAccountDetailActions(props: { account: Account; actionCatalog?: GptActionCatalog; busy: boolean; runWorkflow: AccountWorkflowRunner }) {
  const availability = useGoPayAvailability();
  const actions = goPayActionButtons(props.actionCatalog, props.account, props.busy, availability, props.runWorkflow, 'account_detail');
  return hasVisibleAction(actions) ? <AccountActionRow label="激活" actions={actions} /> : null;
}

export function GoPayAccountRowActions(props: { account: Account; actionCatalog?: GptActionCatalog; busy: boolean; runWorkflow: AccountWorkflowRunner }) {
  const availability = useGoPayAvailability();
  const actions = goPayRowActions(props.actionCatalog, props.account, props.busy, availability, props.runWorkflow, 'account_row');
  return <AccountRowActionGroups actions={actions} />;
}

export function GoPayManualPaymentActions(props: WorkflowJobActionRendererProps) {
  const payment = manualGoPayPaymentView(props.job);
  const canConfirm = canConfirmManualGoPayPayment(props.job, props.progress, payment);
  const runner = useAsyncActionRunner();
  if (!canConfirm || !payment) return null;
  return <ManualPaymentCard payment={payment} busy={runner.isActive('payment')} onConfirm={() => void confirmPayment(props, runner)} />;
}

function goPayActionButtons(catalog: GptActionCatalog | undefined, account: Account, busy: boolean, availability: GoPayAvailability, run: AccountWorkflowRunner, placement: GptActionPlacement): ActionButtonDescriptor[] {
  return privatePaymentActions(catalog, placement).map((action) => {
    const state = accountActionAvailability(catalog, action.action_id, account, placement);
    const allowed = canGoPayPayment(account);
    return {
      id: actionButtonID(action, placement),
      visible: state.visible,
      label: accountActionLabel(catalog, action.action_id, action.display_name || action.action_id, placement),
      hint: !availability.available ? availability.reason : state.reason || (allowed ? 'GoPay 激活支付渠道' : '需要已注册且未激活账号'),
      icon: <span className="activationPaymentIcon"><WalletCards size={15} /></span>,
      className: 'activationAction',
      disabled: busy || !availability.available || !state.enabled || !allowed,
      onClick: () => void run(action.action_id, account, goPayPayload(account), placement, action.display_name),
    };
  });
}

function goPayRowActions(catalog: GptActionCatalog | undefined, account: Account, busy: boolean, availability: GoPayAvailability, run: AccountWorkflowRunner, placement: GptActionPlacement): RowActionDescriptor[] {
  return goPayActionButtons(catalog, account, busy, availability, run, placement)
    .filter((action) => action.visible !== false)
    .map((action) => ({ id: action.id, label: action.label, icon: action.icon || <WalletCards size={14} />, disabled: action.disabled, className: action.className, onClick: () => action.onClick?.() }));
}

function privatePaymentActions(catalog: GptActionCatalog | undefined, placement: string) {
  return (catalog?.actions || []).filter((action) => action.owner === PRIVATE_OWNER && action.capabilities.includes(PAYMENT_CAPABILITY) && action.ui_buttons.some((button) => button.placement === placement));
}

function canGoPayPayment(account: Account) {
  return !isInvalidGptAccount(account) && !isUserAlreadyExistsAccount(account) && !accountIsActivated(account) && hasRegisteredSession(account);
}

function goPayPayload(account: Account) {
  const accountID = accountCarrierID(account);
  return { gopay_account_id: accountID, tokenization: 'true' };
}

function actionButtonID(action: { action_id: string; ui_buttons: Array<{ id: string; placement: string }> }, placement: string) {
  return action.ui_buttons.find((button) => button.placement === placement)?.id || action.action_id;
}

function useGoPayAvailability(): GoPayAvailability {
  const query = useQuery({ queryKey: ['gopay', 'availability'], queryFn: () => api<{ success?: boolean; ok?: boolean }>('/api/gopay/health'), retry: false, staleTime: 15000, refetchInterval: 30000 });
  const available = Boolean(query.data?.success || query.data?.ok);
  return { available, loading: query.isLoading, reason: available ? '' : (query.isLoading ? '检查 gopay-app...' : 'gopay-app 不可用') };
}

async function confirmPayment(props: WorkflowJobActionRendererProps, runner: ReturnType<typeof useAsyncActionRunner>) {
  await runner.tryRun('payment', async () => {
    const resp = await api<{ success?: boolean; error_message?: string }>(`/api/gpt/jobs/${props.job.job_id}/gopay-payment/confirm`, { method: 'POST', body: '{}' });
    props.onMessage?.(resp.error_message ? 'error' : 'ok', resp.error_message || '已确认支付，继续后续步骤');
    if (!resp.error_message) await props.onChanged?.();
  }, { onError: props.onError });
}

function ManualPaymentCard({ payment, busy, onConfirm }: { payment: NonNullable<ReturnType<typeof manualGoPayPaymentView>>; busy: boolean; onConfirm: () => void }) {
  const dataUrl = useQRCodeDataURL(payment.qr_payload);
  return (
    <div className="gptWorkflowActionCard">
      <div className="gptWorkflowActionHead"><strong>扫码支付 QRIS</strong><span>{payment.charge_ref || '等待人工确认'}</span></div>
      <div className="gptWorkflowQRIS">
        {dataUrl ? <img src={dataUrl} alt="QRIS" /> : <div className="gptWorkflowQRPlaceholder">QRIS</div>}
        <div className="gptWorkflowQRISInfo">
          {payment.qr_url && <a href={payment.qr_url} target="_blank" rel="noreferrer">打开远端码<ExternalLink size={12} /></a>}
          {payment.deeplink_url && <a href={payment.deeplink_url} target="_blank" rel="noreferrer">GoPay deeplink<ExternalLink size={12} /></a>}
          <Button type="button" disabled={busy} onClick={onConfirm} {...buttonHint('确认已完成 QRIS 支付')}>{busy ? '确认中' : '我已支付，继续'}</Button>
        </div>
      </div>
    </div>
  );
}

function manualGoPayPaymentView(job: Job) {
  const data = jobDataObject((job.steps || []).find((item) => item.step_name === 'gopay_payment')?.detail);
  const confirmation = objectValue(data.manual_payment_confirmation);
  const complete = objectValue(data.payment_complete);
  const required = confirmation.required === true || complete.awaiting_manual_confirmation === true;
  if (!required) return null;
  const qrValue = stringValue(complete.qr_string) || stringValue(data.qr_string) || stringValue(complete.qr_code_url);
  const qrCodeUrl = stringValue(complete.qr_code_url);
  return { confirmed: confirmation.confirmed === true, charge_ref: stringValue(complete.charge_ref), qr_payload: isHttpURL(qrValue) ? '' : qrValue, qr_url: isHttpURL(qrCodeUrl) ? qrCodeUrl : '', deeplink_url: stringValue(complete.deeplink_url) };
}

function canConfirmManualGoPayPayment(job: Job, progress: WorkflowProgress | null, payment: ReturnType<typeof manualGoPayPaymentView>) {
  return !!payment && !payment.confirmed && job.status === 'RUNNING' && (progress?.step_name === 'gopay_payment' || job.last_step === 'gopay_payment');
}

function useQRCodeDataURL(payload: string) {
  const [dataUrl, setDataUrl] = useState('');
  useEffect(() => {
    let alive = true;
    if (!payload) { setDataUrl(''); return; }
    QRCode.toDataURL(payload, { errorCorrectionLevel: 'M', margin: 1, width: 192 }).then((value) => { if (alive) setDataUrl(value); }).catch(() => { if (alive) setDataUrl(''); });
    return () => { alive = false; };
  }, [payload]);
  return dataUrl;
}

function isHttpURL(value: string) {
  return /^https?:\/\//i.test(String(value || '').trim());
}
