import { Card, CardTitle, CardDescription, CardHeader } from './card'

export default function Cards({ icon, header, value }: { icon: any; header: any; value: any; }) {
  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <div className="w-6 h-6 flex-shrink-0">
              <img
                src={icon}
                alt=""
                className="w-full h-full object-contain"
              />
            </div>

            <span className="text-lg text-white font-medium">
              {header}
            </span>
          </CardTitle>
          <CardDescription>
            <div>
              <p className="text-sm text-[#92929B]">{value}</p>
            </div></CardDescription>
        </CardHeader>

      </Card>
    </>
  );
}
