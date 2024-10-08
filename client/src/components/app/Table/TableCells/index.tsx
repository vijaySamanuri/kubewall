import { CONFIG_MAPS_ENDPOINT, CUSTOM_RESOURCES_LIST_ENDPOINT, ENDPOINTS_ENDPOINT, HPA_ENDPOINT, INGRESSES_ENDPOINT, NODES_ENDPOINT, ROLE_BINDINGS_ENDPOINT, SECRETS_ENDPOINT, SERVICES_ENDPOINT } from '@/constants';

import { ClusterDetails } from '@/types';
import { ConditionCell } from './conditionCell';
import { CurrentByDesiredCell } from './currentByDesiredCell';
import { DefaultCell } from './defaultCell';
import { MultiValueCell } from './multiValueCell';
import { NameCell } from './nameCell';
import { Skeleton } from '@/components/ui/skeleton';
import { StatusCell } from './statusCell';
import { TimeCell } from './timeCell';

type TableCellType = {
  type: string;
  value: string;
  namespace: string;
  instanceType: string;
  loading: boolean;
  queryParams?: string;
} & ClusterDetails;

const TableCells = ({
  clusterName,
  configName,
  instanceType,
  loading,
  namespace,
  type,
  value,
  queryParams
}: TableCellType) => {
  if (loading) {
    return <Skeleton className="h-4" />;
  }
  if(value === undefined || value === 'undefined') {
    return <DefaultCell cellValue='' />;
  }
  if (type === 'Conditions') {
    return <ConditionCell cellValue={value} />;
  }
  if (type === 'Age' || type === 'Duration' || ((type === 'eventTime' || type === 'firstTimestamp' || type === 'lastTimestamp') && value !== 'null')) {
    return <TimeCell cellValue={value} />;
  }
  if (type === 'Ready' || type === 'Current') {
    return <CurrentByDesiredCell cellValue={value} />;
  }
  if (type === 'Status' || type === 'reason' || type === 'Condition Status') {
    return <StatusCell cellValue={value} />;
  }
  if (type === 'Name') {
    let link = '';
    if (instanceType !== CUSTOM_RESOURCES_LIST_ENDPOINT) {
      link = `${configName}/${clusterName}/details?resourcekind=${instanceType.toLowerCase()}&resourcename=${value}${namespace ? `&namespace=${namespace}` : ''}`;
    } else {
      link = `${configName}/${clusterName}/details?resourcekind=${instanceType.toLowerCase()}&resourcename=${value}&${queryParams}${namespace ? `&namespace=${namespace}` : ''}`;
    }
    return <NameCell
      cellValue={value}
      link={link}
    />;
  }
  if (instanceType === 'events' || instanceType === HPA_ENDPOINT) {
    const eventsValue = value === 'null' ? '—' : value;
    return <DefaultCell cellValue={eventsValue} truncate={false} />;
  }
  if (
    value !== '' &&
    (type === 'Rules' || type === 'Ports' || type === 'Bindings' || type === 'Roles' || type === 'Keys') &&
    (
      instanceType === INGRESSES_ENDPOINT ||
      instanceType === ENDPOINTS_ENDPOINT ||
      instanceType === SERVICES_ENDPOINT ||
      instanceType === ROLE_BINDINGS_ENDPOINT ||
      instanceType === NODES_ENDPOINT || 
      instanceType === SECRETS_ENDPOINT ||
      instanceType === CONFIG_MAPS_ENDPOINT
    )
  ) {
    return <MultiValueCell cellValue={value} />;
  }
  return <DefaultCell cellValue={value === '' ? '—' : value} />;
};

export {
  TableCells
};
