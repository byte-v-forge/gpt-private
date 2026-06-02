import { registerAccountDetailActionExtensions, registerAccountRowActionExtensions } from './account-extension-registry';
import { registerWorkflowJobActionRenderers } from './job-action-renderers';
import { GoPayAccountDetailActions, GoPayAccountRowActions, GoPayManualPaymentActions } from './gopay-dashboard-extension';
import './gopay-workflow.css';

let registered = false;

export function registerPrivateDashboardExtensions() {
  if (registered) return;
  registered = true;
  registerAccountDetailActionExtensions([{ id: 'gpt-private.gopay.detail', Component: GoPayAccountDetailActions }]);
  registerAccountRowActionExtensions([{ id: 'gpt-private.gopay.row', Component: GoPayAccountRowActions }]);
  registerWorkflowJobActionRenderers([{
    id: 'gpt-private.gopay.manual-payment',
    jobActions: ['GOPAY_PAYMENT', 'GOPAY_QRIS_PAYMENT_ACTIVATE'],
    statuses: ['RUNNING'],
    render: (props) => <GoPayManualPaymentActions {...props} />,
  }]);
}
