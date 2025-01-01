import { WebSocketProvider } from "@/components/ui/WebSocketProvider";
import CheckAuth from "@/lib/checkAuth";
import { children } from "solid-js";

export function Layout(props) {
  return (
    <>
      <CheckAuth />
      <WebSocketProvider />
      {props.children}
    </>
  );
}
