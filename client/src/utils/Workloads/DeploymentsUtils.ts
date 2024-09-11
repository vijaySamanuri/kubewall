import { DeploymentsResponse } from "@/types";

const formatDeploymentsResponse = (deployments: DeploymentsResponse[]) => {
  return deployments.map(({namespace, name, status, age, hasUpdated}) => ({
    namespace:namespace,
    name: name,
    ready: `${status.readyReplicas}/${status.replicas}`,
    desired: status.replicas,
    updated: status.updatedReplicas,
    available: status.availableReplicas,
    age: age,
    hasUpdated: hasUpdated,
    conditions: status.conditions.map(({type}) => type).toString()
  }));
};

export {
  formatDeploymentsResponse
};
