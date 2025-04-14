"use client"

import * as React from "react"
import { ChevronDown, ChevronLeft, ChevronRight, ChevronUp } from "lucide-react"
import { DayPicker, UI, useDayPicker } from "react-day-picker"

import { cn } from "@/lib/utils"
import { buttonVariants } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export type CalendarProps = React.ComponentProps<typeof DayPicker>

function Calendar({
  captionLayout = "label",
  className,
  classNames,
  showOutsideDays = true,
  ...props
}: CalendarProps) {
  return (
    <DayPicker
      showOutsideDays={showOutsideDays}
      captionLayout={captionLayout}
      className={cn("p-3", className)}
      classNames={{
        button_next: cn(
          buttonVariants({ variant: "outline" }),
          "size-7 bg-transparent p-0 opacity-50 hover:opacity-100"
        ),
        button_previous: cn(
          buttonVariants({ variant: "outline" }),
          "size-7 bg-transparent p-0 opacity-50 hover:opacity-100"
        ),
        caption_label: "text-sm font-medium aria-hidden:hidden",
        day_button: cn(
          buttonVariants({ variant: "ghost" }),
          "size-8 font-normal"
        ),
        day: "rounded-md p-0 text-center text-sm aria-selected:bg-accent",
        disabled: "*:text-muted-foreground *:opacity-50",
        dropdown: "first:basis-3/5 last:basis-2/5",
        dropdowns: "flex basis-full items-center gap-2 text-sm font-medium",
        hidden: "invisible",
        month_caption: "flex items-center justify-center pt-1",
        month_grid: "w-full border-collapse space-y-1",
        month: cn(
          "space-y-4",
          captionLayout !== "label" && !props.hideNavigation && "mt-9"
        ),
        months:
          "relative flex flex-col gap-y-4 sm:flex-row sm:gap-x-4 sm:gap-y-0",
        nav: "absolute flex w-full items-center justify-between space-x-1 px-1",
        outside:
          "*:text-muted-foreground *:opacity-50 *:aria-selected:bg-accent/50 *:aria-selected:text-muted-foreground *:aria-selected:opacity-30",
        range_end: "rounded-l-none",
        range_middle:
          "rounded-none first:rounded-l-md last:rounded-r-md *:aria-selected:bg-accent *:aria-selected:text-accent-foreground",
        range_start: "rounded-r-none",
        selected:
          "*:bg-primary *:text-primary-foreground *:hover:bg-primary *:hover:text-primary-foreground *:focus:bg-primary *:focus:text-primary-foreground",
        today: "*:bg-accent *:text-accent-foreground",
        week: "mt-2 flex w-full",
        weekday: "w-8 text-[0.8rem] font-normal text-muted-foreground",
        weekdays: "flex",
        ...classNames,
      }}
      components={{
        Chevron: ({ orientation }) => {
          switch (orientation) {
            case "up":
              return <ChevronUp className="size-4" />
            case "down":
              return <ChevronDown className="size-4" />
            case "left":
              return <ChevronLeft className="size-4" />
            case "right":
            default:
              return <ChevronRight className="size-4" />
          }
        },
        Dropdown: ({
          "aria-label": ariaLabel,
          disabled,
          value,
          onChange,
          options,
          className,
        }) => {
          const { classNames } = useDayPicker()

          return (
            <Select
              disabled={disabled}
              value={`${value}`}
              onValueChange={(value) =>
                onChange?.({
                  target: { value },
                } as React.ChangeEvent<HTMLSelectElement>)
              }
            >
              <SelectTrigger
                aria-label={ariaLabel}
                className={cn(classNames[UI.Dropdown], className)}
              >
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {options?.map((option) => (
                  <SelectItem
                    key={option.value}
                    value={`${option.value}`}
                    disabled={option.disabled}
                  >
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )
        },
      }}
      {...props}
    />
  )
}
Calendar.displayName = "Calendar"

export { Calendar }

