import heroimg from "@/assets/heroimg.svg";
import { Link } from "react-router";
import { Button } from "./button";
import { WavyBackground } from "./wavy-background";
import { Github } from 'lucide-react';
import CardsHero from "./heroCards";
const hero = () => {
  return (
    <>
      <WavyBackground>

        <div className="lg:flex md:grid-cols-1  ">
          <div className=" pt-32 pl-32  ">
            <p className=" lg:text-5xl md:text-lg sm:text-xs text-white font-bold">
              OpenFin<span className="text-pink-500"> elevate </span>
              your financial
              <br /> journey with empowering tools
              <br /> and insightful guidance
            </p>

            <p className=" lg:text-lg md:text-sm sm:text-xs text-[#95959D]  pt-4 pb-4">
              Empower your finances with OpenFin's smart tools and expert guidance
            </p>
            <div className="space-x-5 flex">
              <Link to="/auth">
                <Button
                  variant="link"
                >
                  <span>Get Started</span>
                  <svg
                    className="inline"
                    width="20"
                    height="10"
                    viewBox="0 0 20 10"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M15.0125 3.75H0V6.25H15.0125V10L20 5L15.0125 0V3.75Z"
                      fill="white"
                    />
                  </svg>
                </Button>
              </Link>
              <Link to="https://github.com/xDeFc0nx/OpenFin">
                <Button
                >
                  <Github />
                  Github
                </Button>
              </Link>
            </div>
          </div>
          <div className="pl-56 pt-20 flex items-center ">
            <div
              className={` flex text-white w-[30rem] h-[20rem] bg-white/10  backdrop-filter backdrop-blur-md rounded-lg 	   shadow-lg`}
            >
              <img
                className="w-full rounded-md"
                src={heroimg}
                alt="financial svg"
              />
            </div>
          </div>
        </div>
          <CardsHero/>
      </WavyBackground>
    </>
  );
};

export default hero;
