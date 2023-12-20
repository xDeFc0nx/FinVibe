/* eslint-disable import/no-extraneous-dependencies */
/* eslint no-underscore-dangle: 0 */

import React from "react";
import { Card, CardHeader, CardBody } from "@nextui-org/card";
import Image from "next/image";

export default function card({ icon, header, value }) {
  return (
    <Card className="w-[20rem] h-[13rem]  text-white px-2 py-2  bg-white/10  backdrop-filter backdrop-blur-lg	 shadow-lg">
      <CardHeader className="w-full  justify-start">
        <Image height={30} width={30} radius="sm" src={icon} />

        <div className="flex flex-col justify-center">
          <p className="text-lg  text-white pl-1">{header}</p>
        </div>
      </CardHeader>

      <CardBody>
        <p className="text-sm text-[#92929B]">{value}</p>
      </CardBody>
    </Card>
  );
}

export function CardDashboard({ icon, header, value, color }) {
  return (
    <div className=" col-span-1 p-10 text-white  bg-secondary-gray/50  backdrop-filter backdrop-blur-lg	rounded-lg  shadow-lg">
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
