import { objectValue, stringValue } from '@byte-v-forge/common-ui';
import { stepDetailData } from './job-utils';
import type { ConcreteGoPayPaymentChannel, Job, WorkflowProgress } from './types';

type PaymentChannelDescriptor = {
  value: ConcreteGoPayPaymentChannel;
  canonical: '' | 'gopay_sms' | 'gopay_wa';
  label: string;
  aliases: string[];
  match?: (value: string) => boolean;
};

const PAYMENT_CHANNELS: PaymentChannelDescriptor[] = [{
  value: 'account',
  canonical: '',
  label: 'GoPay 支付',
  aliases: ['account', 'gopay', 'gopay_account']
}, {
  value: 'sms',
  canonical: 'gopay_sms',
  label: 'Gopay-SMS',
  aliases: ['sms', 'gopay_sms', 'gopay-sms'],
  match: (value) => value.includes('sms') && value.includes('gopay')
}, {
  value: 'app_wa',
  canonical: '',
  label: 'Gopay App-WA',
  aliases: ['app_wa', 'gopay_app_wa', 'gopay-app-wa']
}, {
  value: 'wa',
  canonical: 'gopay_wa',
  label: 'Gopay-WA',
  aliases: ['wa', 'whatsapp', 'gopay_wa', 'gopay-wa'],
  match: (value) => (value.includes('wa') || value.includes('whatsapp')) && value.includes('gopay')
}];

export const GO_PAY_PAYMENT_CHANNELS: ConcreteGoPayPaymentChannel[] = ['account'];

export function paymentChannelValue(value: string): '' | 'gopay_sms' | 'gopay_wa' {
  const normalized = String(value || '').trim().toLowerCase();
  if (!normalized) return '';
  return paymentChannelDescriptor(normalized)?.canonical || '';
}

export function goPayPaymentChannelLabel(value: string) {
  const normalized = String(value || '').trim().toLowerCase();
  return paymentChannelDescriptor(normalized)?.label || '-';
}

export function goPayPaymentActionLabel(channel: ConcreteGoPayPaymentChannel) {
  return paymentChannelDescriptor(channel)?.label || '-';
}

export function manualGoPayPaymentView(job: Job) {
  const data = stepDetailData((job.steps || []).find((item) => item.step_name === 'gopay_payment'));
  if (!data) return null;
  const confirmation = objectValue(data.manual_payment_confirmation);
  const complete = objectValue(data.payment_complete);
  const required = confirmation.required === true || complete.awaiting_manual_confirmation === true;
  if (!required) return null;
  const qrValue = stringValue(complete.qr_string) || stringValue(data.qr_string) || stringValue(complete.qr_code_url);
  const qrCodeUrl = stringValue(complete.qr_code_url);
  return {
    required,
    confirmed: confirmation.confirmed === true,
    charge_ref: stringValue(complete.charge_ref),
    qr_payload: isHttpURL(qrValue) ? '' : qrValue,
    qr_url: isHttpURL(qrCodeUrl) ? qrCodeUrl : '',
    deeplink_url: stringValue(complete.deeplink_url)
  };
}

export function canConfirmManualGoPayPayment(job: Job, progress: WorkflowProgress | null, payment: ReturnType<typeof manualGoPayPaymentView>) {
  return !!payment && !payment.confirmed && job.status === 'RUNNING' &&
    (progress?.step_name === 'gopay_payment' || job.last_step === 'gopay_payment');
}

function isHttpURL(value: string) {
  return /^https?:\/\//i.test(String(value || '').trim());
}

function paymentChannelDescriptor(value: string) {
  const normalized = String(value || '').trim().toLowerCase();
  return PAYMENT_CHANNELS.find((item) => item.aliases.includes(normalized) || item.match?.(normalized));
}
