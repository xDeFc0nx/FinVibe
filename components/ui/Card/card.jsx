/* eslint-disable import/no-extraneous-dependencies */
/* eslint no-underscore-dangle: 0 */

import React from "react";
import Image from "next/image";

export default function card({ icon, header, value }) {
  return (
    <div className="w-[20rem] h-[13rem]  text-white px-4 py-4 rounded-md bg-white/10  backdrop-filter backdrop-blur-lg	 shadow-lg">
      <div className="w-full flex justify-start">
        <Image height={30} width={30} radius="sm" src={icon} />

        <div className="flex flex-col justify-center">
          <p className="text-lg  text-white pl-1">{header}</p>
        </div>
      </div>

      <div>
        <p className="text-sm text-[#92929B]">{value}</p>
      </div>
    </div>
  );
}

export function CardDashboard({ icon, header, value, color }) {
  return (
    <div className=" col-span-1 p-10 text-white  bg-secondary-gray/50  backdrop-filter backdrop-blur-lg	rounded-lg  shadow-inner">
      <div className="w-full flex justify-start">
        <Image height={20} width={20} radius="sm" src={icon} />

        <div className="flex flex-col justify-center">
          <p className="text-lg  text-[#92929B] pl-1">{header}</p>
        </div>
      </div>

      <div className="pt-4">
        <p className={`text-4xl ${color}  `}>{value}</p>
      </div>
    </div>
  );
}
export function CardChart({ header, styles, children }) {
  return (
    <div
      className={`"${styles}col-span-1 p-10 text-white bg-secondary-gray/50 backdrop-filter backdrop-blur-lg rounded-lg shadow-inner"`}
    >
      <div className="flex justify-start">
        <div className="flex flex-col justify-center">
          <p className="text-lg text-[#92929B] pl-1">{header}</p>
          <div className="pt-4">{children}</div>
        </div>
      </div>
    </div>
  );
}
