import { WebSocketProvider } from "@/components/ui/WebSocketProvider";
import CheckAuth from "@/lib/checkAuth";

export function Layout(props) {
  return (
    <>
      <CheckAuth />
      <WebSocketProvider />
      {props.children}
    </>
  );
}
