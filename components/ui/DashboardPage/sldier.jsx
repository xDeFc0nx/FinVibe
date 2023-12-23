import React from "react";

function Slider({ active, text }) {
  return (
    <div>
      <button
        type="button"
        className={`
    ${active ? " bg-green-500" : "text-white "}
    h-10 w-32 rounded-full`}
      >
        {text}
      </button>
    </div>
  );
}

export default Slider;
