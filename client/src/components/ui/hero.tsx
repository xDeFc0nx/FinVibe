import heroimg from "@/assets/heroimg.svg";
import { Link } from "react-router-dom";
import { Button } from "./button"; // Assuming Button component is responsive
import { WavyBackground } from "./wavy-background"; // Assuming WavyBackground handles its own responsiveness or is purely decorative
import { Github } from 'lucide-react';
import CardsHero from "./heroCards"; // Assuming CardsHero is responsive

const Hero = () => {
  return (
    <>
      <WavyBackground>

        <div className="max-w-screen-xl mx-auto px-4 py-16 sm:px-6 lg:px-8 lg:py-24">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12 lg:gap-16 md:items-center">

            <div>
              <h1 className="text-3xl font-bold text-white sm:text-4xl md:text-5xl lg:text-5xl">
                NovaoFin<span className="text-pink-500"> elevate </span>
                your financial
                journey with empowering tools
                <br /> and insightful guidance
              </h1>

              <p className="mt-4 text-lg text-[#95959D] md:text-xl">
                Empower your finances with NovaoFin's smart tools and expert guidance
              </p>

              <div className="mt-8 flex flex-col space-y-4 sm:flex-row sm:space-y-0 sm:space-x-4 items-start sm:items-center">
                <Link to="/auth">
                  <Button
                    variant="link"
                    className="inline-flex items-center text-base font-medium"
                  >
                    <span>Get Started</span>
                    <svg
                      className="ml-2"
                      width="20"
                      height="10"
                      viewBox="0 0 20 10"
                      fill="currentColor"
                      xmlns="http://www.w3.org/2000/svg"
                      aria-hidden="true"
                    >
                      <path
                        d="M15.0125 3.75H0V6.25H15.0125V10L20 5L15.0125 0V3.75Z"
                      />
                    </svg>
                  </Button>
                </Link>
                <Link to="https://github.com/xDeFc0nx/NovaoFin" target="_blank" rel="noopener noreferrer">
                  <Button
                    className="inline-flex items-center text-base font-medium"
                  >
                    <Github className="mr-2 h-5 w-5" />
                    Github
                  </Button>
                </Link>
              </div>
            </div>

            <div className="mt-8 md:mt-0 mx-auto md:mx-0 w-full max-w-md md:max-w-full">
              <div className={`relative text-white bg-white/10 backdrop-filter backdrop-blur-md rounded-lg shadow-lg overflow-hidden`}>
                <img
                  className="w-full h-auto object-cover rounded-md"
                  src={heroimg}
                  alt="Illustration of financial tools and charts"
                />
              </div>
            </div>

          </div>

          <div className="mt-16 sm:mt-24 lg:mt-32"> {/* Replaced pt-24 with responsive mt-* */}
            <CardsHero />
          </div>
        </div>
      </WavyBackground>
    </>
  );
};

export default Hero;
