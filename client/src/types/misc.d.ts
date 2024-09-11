type HeaderList = {
  title: string;
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  accessorFn?: (property: any) => void
  accessorKey?: string;
  enableSorting?: boolean;
}

type ResourceDetailsSearch = {
  namespace: string;
};

type kwListSearch = {
  resourcekind: string;
  group?: string;
  kind?: string;
  resource?: string;
  version?: string;
};

type kwDetailsSearch = {
  resourcename: string;
  resourcekind: string;
  group?: string;
  kind?: string;
  resource?: string;
  version?: string;
  namespace?: string;
};

type CustomResoucesListSearch = {
  group: string;
  version: string;
  resource: string;
  kind: string;
};

type NavigationRoutes = {
  [key: string]: {
    route: string,
    name: string
  }[];
};

type ClusterDetails = {
  configName: string;
  clusterName: string;
};

type TableColumns = {
  count: number;
  configName: string;
  clusterName: string;
  instanceType: string;
  loading: boolean;
  queryParams?: string;
};

type CustomResources = {
  activeVersion: string;
  age: string,
  hasUpdated: boolean,
  name: string,
  scope: string,
  queryParam: string,
  spec: {
      conversion: {
          strategy: string
      },
      group: string,
      names: {
          kind: string,
          listKind: string,
          plural: string,
          shortNames: string | null,
          singular: string
      },
      scope: string
  },
  versions: number
};

type CustomResourcesNavigationKeys = {
  resources: {
    name: string;
    route: string;
  }[];
};

type CustomResourcesNavigation = {
  [key: string]: CustomResourcesNavigationKeys;
};

type KeyValue = {
  [key: string]: string
};

type KeyValueNull = ({
  [key: string]:  string | number | null
}| null);

export {
  ClusterDetails,
  CustomResources,
  CustomResoucesListSearch,
  CustomResourcesNavigationKeys,
  CustomResourcesNavigation,
  HeaderList,
  ResourceDetailsSearch,
  NavigationRoutes,
  TableColumns,
  kwListSearch,
  kwDetailsSearch,
  KeyValue,
  KeyValueNull
};