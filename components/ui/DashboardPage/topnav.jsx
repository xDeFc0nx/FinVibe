"use client";

import React, { useState } from "react";
import Link from "next/link";
import Button from "../Button";
import Slider from "./sldier";

function TopNav() {
  const [activeIndex, setActiveIndex] = useState(1);

  return (
    <div>
      <div className="flex gap-2 items-center ">
        <Link href="/dashboard/new">
          <Button Text="Add New" color="bg-green-500" type="Button" />
        </Link>
        <div className="bg-secondary-gray/50 backdrop-filter backdrop-blur-lg shadow-inner mb-2 flex rounded-full">
          <Slider
            text="Previous month"
            active={activeIndex === 0}
            onClick={() => setActiveIndex(0)}
          />
          <Slider
            text="Current month"
            active={activeIndex === 1}
            onClick={() => setActiveIndex(1)}
          />
          <Slider
            text="All time"
            active={activeIndex === 2}
            onClick={() => setActiveIndex(2)}
          />
        </div>
      </div>
    </div>
  );
}

export default TopNav;
