import { Component, JSX, Suspense } from "solid-js";

const App: Component<{ children: JSX.Element }> = (props) => {
  return (
    <>
      <Suspense>{props.children}</Suspense>
    </>
  );
};

export default App;
