import { registerWorkflowStepRenderers } from './step-renderers';

let registered = false;

export function registerGoPayWorkflowRenderers() {
  if (registered) return;
  registered = true;
  registerWorkflowStepRenderers([
    {
      id: 'gpt.gopay.pin',
      stepNames: ['gopay_app_ensure_pin_setup'],
      label: 'GoPay 确认 PIN'
    }
  ]);
}
