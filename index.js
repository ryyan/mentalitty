const reblessed = require("reblessed");

// Create a screen object.
const screen = reblessed.screen({
  smartCSR: true,
});

screen.title = "MentaliTTY";

// Create a box perfectly centered horizontally and vertically.
const box = reblessed.box({
  top: "center",
  left: "center",
  width: "50%",
  height: "50%",
  content: "Hello {bold}world{/bold}!",
  tags: true,
  border: {
    type: "line",
  },
  style: {
    fg: "white",
    bg: "magenta",
    border: {
      fg: "#f0f0f0",
    },
    hover: {
      bg: "green",
    },
  },
});

// Append our box to the screen.
screen.append(box);

// If box is focused, handle `enter`/`return` and give us some more content.
box.key("enter", (ch, key) => {
  box.setContent(
    "{right}Even different {black-fg}content{/black-fg}.{/right}\n"
  );
  box.setLine(1, "bar");
  box.insertLine(1, "foo");
  screen.render();
});

// Quit on Escape, q, or Control-C.
screen.key(["escape", "q", "C-c"], (ch, key) => {
  return process.exit(0);
});

// Focus our element.
box.focus();

// Render the screen.
screen.render();
