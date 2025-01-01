import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import { Layout } from "@/pages/dashboard/layout";
import index from "@/pages/index";
import dashboard from "@/pages/dashboard/index";
import Login from "@/pages/login";
import { Toaster } from "solid-toast";

const App = (props) => (
  <>
    {" "}
    <Toaster />
    {props.children}
  </>
);

render(
  () => (
    <Router root={App}>
      <Route path="/login" component={Login} />
      <Route path="/app" component={Layout}>
        <Route path="/dashboard" component={dashboard} />
        <Route path="/" component={index} />
      </Route>
    </Router>
  ),
  document.getElementById("root")
);

export default App;
