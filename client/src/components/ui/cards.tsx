

export default function Cards({ icon, header, value }: { icon: any; header: any; value: any; }) {
  return (
    <div className="w-[20rem] h-[13rem]  text-white px-4 py-4 rounded-md bg-white/10  backdrop-filter backdrop-blur-lg	 shadow-lg">
      <div className="w-full flex justify-start">
        <img src={icon} />

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
