import { render, screen } from "@testing-library/solid";
import { describe, it, expect } from "vitest";

function App() {
  return <h1>openCenter-base</h1>;
}

describe("App", () => {
  it("renders a headline", () => {
    render(() => <App />);
    expect(screen.getByRole("heading", { name: "openCenter-base" })).toBeInTheDocument();
  });
});
