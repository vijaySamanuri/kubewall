import { Outlet, useRouterState } from "@tanstack/react-router";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { createEventStreamQueryObject, getEventStreamUrl } from "./utils";

import { NAMESPACES_ENDPOINT } from "./constants";
import { NamespacesResponse } from "./types";
import { Sidebar } from "@/components/app/Sidebar";
import { updateNamspaces } from "./data/Clusters/Namespaces/NamespacesSlice";
import { useAppDispatch } from "./redux/hooks";
import { useEventSource } from "./components/app/Common/Hooks/EventSource";

export function App() {
  const dispatch = useAppDispatch();
  const router = useRouterState();
  const pathname = router.location.pathname;
  const configName = pathname.split('/')[1];
  const clusterName = pathname.split('/')[2];


  const sendMessage = (message: NamespacesResponse[]) => {
    dispatch(updateNamspaces(message));
  };

  useEventSource({
    url: getEventStreamUrl(
      NAMESPACES_ENDPOINT,
      createEventStreamQueryObject(
        configName,
        clusterName
      )),
    sendMessage
  });

return (
  <ResizablePanelGroup
    direction="horizontal"
    className="h-screen items-stretch"
  >
    <ResizablePanel defaultSize={11}>
      <Sidebar />
    </ResizablePanel>
    <ResizableHandle withHandle />
    <ResizablePanel defaultSize={89}>
      <Outlet />
    </ResizablePanel>
  </ResizablePanelGroup>
);
}