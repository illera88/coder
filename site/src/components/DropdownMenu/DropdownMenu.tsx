/**
 * Copied from shadc/ui on 12/19/2024
 * @see {@link https://ui.shadcn.com/docs/components/dropdown-menu}
 *
 * This component was updated to match the styles from the Figma design:
 * @see {@link https://www.figma.com/design/WfqIgsTFXN2BscBSSyXWF8/Coder-kit?node-id=656-2354&t=CiGt5le3yJEwMH4M-0}
 */

import * as DropdownMenuPrimitive from "@radix-ui/react-dropdown-menu";
import { Circle } from "lucide-react";
import type * as React from "react";
import { cn } from "utils/cn";

export const DropdownMenu = DropdownMenuPrimitive.Root;

export const DropdownMenuTrigger = DropdownMenuPrimitive.Trigger;

export const DropdownMenuGroup = DropdownMenuPrimitive.Group;

export const DropdownMenuRadioGroup = DropdownMenuPrimitive.RadioGroup;

export const DropdownMenuContent: React.FC<
	React.ComponentPropsWithRef<typeof DropdownMenuPrimitive.Content>
> = ({ className, sideOffset = 4, ...props }) => (
	<DropdownMenuPrimitive.Portal>
		<DropdownMenuPrimitive.Content
			sideOffset={sideOffset}
			className={cn(
				[
					"z-50 min-w-48 overflow-hidden rounded-md border border-solid bg-surface-primary p-2 text-content-secondary shadow-md",
					"data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
					"data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
					"data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2",
					"data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2",
				],
				className,
			)}
			{...props}
		/>
	</DropdownMenuPrimitive.Portal>
);

type DropdownMenuItemProps = React.ComponentPropsWithRef<
	typeof DropdownMenuPrimitive.Item
> & {
	inset?: boolean;
};

export const DropdownMenuItem: React.FC<DropdownMenuItemProps> = ({
	children,
	className,
	inset,
	...props
}) => (
	<DropdownMenuPrimitive.Item
		className={cn(
			`
			relative flex cursor-default gap-2 select-none items-center px-3 py-2.5 text-sm font-medium outline-none transition-colors
			hover:bg-surface-secondary text-content-secondary hover:text-content-primary
			data-[disabled]:pointer-events-none data-[disabled]:opacity-50 data-[disabled]:text-content-disabled`,
			inset && "pl-8",
			className,
		)}
		{...props}
	>
		{children}
	</DropdownMenuPrimitive.Item>
);

export const DropdownMenuRadioItem: React.FC<
	React.ComponentPropsWithRef<typeof DropdownMenuPrimitive.RadioItem>
> = ({ className, children, ...props }) => {
	return (
		<DropdownMenuPrimitive.RadioItem
			className={cn(
				[
					"relative flex cursor-default select-none items-center rounded-sm py-1.5 pr-8 pl-2 text-sm outline-none transition-colors",
					"focus:bg-surface-secondary focus:text-content-primary data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
					"data-[state=checked]:bg-surface-secondary data-[state=checked]:text-content-primary",
				],
				className,
			)}
			{...props}
		>
			{children}
			<span className="absolute right-2 flex h-3.5 w-3.5 items-center justify-center">
				<DropdownMenuPrimitive.ItemIndicator>
					<Circle className="h-2 w-2 fill-current" />
				</DropdownMenuPrimitive.ItemIndicator>
			</span>
		</DropdownMenuPrimitive.RadioItem>
	);
};

export const DropdownMenuSeparator: React.FC<
	React.ComponentPropsWithRef<typeof DropdownMenuPrimitive.Separator>
> = ({ className, ...props }) => {
	return (
		<DropdownMenuPrimitive.Separator
			className={cn("-mx-1 my-2 h-px bg-border", className)}
			{...props}
		/>
	);
};

export const DropdownMenuShortcut: React.FC<
	React.ComponentPropsWithRef<"span">
> = ({ className, ...props }) => {
	return (
		<span
			className={cn("ml-auto text-xs tracking-widest opacity-60", className)}
			{...props}
		/>
	);
};
