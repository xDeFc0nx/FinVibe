/* eslint-disable import/no-extraneous-dependencies */
/* eslint no-underscore-dangle: 0 */

import React from "react";
import { Card, CardHeader, CardBody } from "@nextui-org/card";
import Image from "next/image";

export default function card({ icon, header, value }) {
  return (
    <Card className="w-[20rem] h-[13rem]  text-white  bg-secondary-gray/20  backdrop-filter backdrop-blur-lg	   shadow-lg">
      <CardHeader className="w-full justify-start">
        <Image height={30} width={30} radius="sm" src={icon} />

        <div className="flex flex-col justify-center">
          <p className="text-lg  text-white ">{header}</p>
        </div>
      </CardHeader>

      <CardBody>
        <p className="text-sm text-[#92929B]">{value}</p>
      </CardBody>
    </Card>
  );
}
