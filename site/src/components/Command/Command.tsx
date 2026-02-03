/**
 * Copied from shadc/ui on 11/13/2024
 * @see {@link https://ui.shadcn.com/docs/components/command}
 */
import type { DialogProps } from "@radix-ui/react-dialog";
import { Command as CommandPrimitive } from "cmdk";
import { Dialog, DialogContent } from "components/Dialog/Dialog";
import { Search } from "lucide-react";
import { type FC, type Ref } from "react";
import { cn } from "utils/cn";

export const Command = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive>) => (
	<CommandPrimitive
		ref={ref}
		className={cn(
			"flex h-full w-full flex-col overflow-hidden rounded-md bg-surface-primary text-content-primary",
			className,
		)}
		{...props}
	/>
);
Command.displayName = CommandPrimitive.displayName;

export const CommandDialog: FC<DialogProps> = ({ children, ...props }) => {
	return (
		<Dialog {...props}>
			<DialogContent className="overflow-hidden p-0">
				<Command className="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-content-secondary [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-group]]:px-2 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-input]]:h-12 [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-item]_svg]:h-5 [&_[cmdk-item]_svg]:w-5">
					{children}
				</Command>
			</DialogContent>
		</Dialog>
	);
};

export const CommandDialog2: FC<DialogProps> = ({ children, ...props }) => {
	return (
		<Dialog {...props}>
			<DialogContent variant="withPadding" className="overflow-hidden">
				<Command className="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-content-secondary [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-group]]:px-2 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-input]]:h-12 [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-item]_svg]:h-5 [&_[cmdk-item]_svg]:w-5">
					{children}
				</Command>
			</DialogContent>
		</Dialog>
	);
};

export const CommandInput = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive.Input>) => (
	<div className="flex items-center border-b px-3" cmdk-input-wrapper="">
		<Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
		<CommandPrimitive.Input
			ref={ref}
			className={cn(
				`flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none border-none
				placeholder:text-content-secondary text-content-primary
				disabled:cursor-not-allowed disabled:opacity-50`,
				className,
			)}
			{...props}
		/>
	</div>
);

CommandInput.displayName = CommandPrimitive.Input.displayName;

export const CommandList = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive.List>) => (
	<CommandPrimitive.List
		ref={ref}
		className={cn(
			"max-h-96 overflow-y-auto overflow-x-hidden border-0 border-t border-solid border-border",
			className,
		)}
		{...props}
	/>
);

CommandList.displayName = CommandPrimitive.List.displayName;

export const CommandEmpty = ({
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive.Empty>) => (
	<CommandPrimitive.Empty
		ref={ref}
		className="py-6 text-center text-sm"
		{...props}
	/>
);

CommandEmpty.displayName = CommandPrimitive.Empty.displayName;

export const CommandGroup = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive.Group>) => (
	<CommandPrimitive.Group
		ref={ref}
		className={cn(
			`overflow-hidden p-2 text-content-primary
			[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1.5 [&_[cmdk-group-heading]]:text-xs
			[&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-content-secondary`,
			className,
		)}
		{...props}
	/>
);

CommandGroup.displayName = CommandPrimitive.Group.displayName;

export const CommandSeparator = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive.Separator>) => (
	<CommandPrimitive.Separator
		ref={ref}
		className={cn("-mx-1 h-px bg-border", className)}
		{...props}
	/>
);
CommandSeparator.displayName = CommandPrimitive.Separator.displayName;

export const CommandItem = ({
	className,
	ref,
	...props
}: React.ComponentPropsWithRef<typeof CommandPrimitive.Item>) => (
	<CommandPrimitive.Item
		ref={ref}
		className={cn(
			`relative flex cursor-default gap-2 select-none text-content-secondary items-center rounded-sm px-2 py-2 text-sm font-medium outline-none
			data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50
			data-[selected=true]:bg-surface-secondary data-[selected=true]:text-content-primary
			[&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0`,
			className,
		)}
		{...props}
	/>
);

CommandItem.displayName = CommandPrimitive.Item.displayName;

export const CommandShortcut = ({
	className,
	...props
}: React.HTMLAttributes<HTMLSpanElement>) => {
	return (
		<span
			className={cn(
				"ml-auto text-xs tracking-widest text-content-tertiary",
				className,
			)}
			{...props}
		/>
	);
};
CommandShortcut.displayName = "CommandShortcut";
