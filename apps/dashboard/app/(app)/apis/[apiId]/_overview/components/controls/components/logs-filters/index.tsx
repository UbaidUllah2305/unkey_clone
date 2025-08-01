import { FiltersPopover } from "@/components/logs/checkbox/filters-popover";

import { FilterOperatorInput } from "@/components/logs/filter-operator-input";
import { BarsFilter } from "@unkey/icons";
import { Button } from "@unkey/ui";
import { cn } from "@unkey/ui/src/lib/utils";
import { useState } from "react";
import { keysOverviewFilterFieldConfig } from "../../../../filters.schema";
import { useFilters } from "../../../../hooks/use-filters";
import { OutcomesFilter } from "./outcome-filter";

export const LogsFilters = () => {
  const { filters, updateFilters } = useFilters();
  const [open, setOpen] = useState(false);

  const options = keysOverviewFilterFieldConfig.names.operators.map((op) => ({
    id: op,
    label: op,
  }));
  const activeNameFilter = filters.find((f) => f.field === "names");
  const activeIdentityFilter = filters.find((f) => f.field === "identities");
  const activeKeyIdsFilter = filters.find((f) => f.field === "keyIds");
  const activeTagsFilter = filters.find((f) => f.field === "tags");
  const keyIdOptions = keysOverviewFilterFieldConfig.keyIds.operators.map((op) => ({
    id: op,
    label: op,
  }));
  const tagsOptions = keysOverviewFilterFieldConfig.tags.operators.map((op) => ({
    id: op,
    label: op,
  }));
  return (
    <FiltersPopover
      open={open}
      onOpenChange={setOpen}
      items={[
        {
          id: "names",
          label: "Name",
          shortcut: "n",
          component: (
            <FilterOperatorInput
              label="Name"
              options={options}
              defaultOption={activeNameFilter?.operator}
              defaultText={activeNameFilter?.value as string}
              onApply={(id, text) => {
                const activeFiltersWithoutNames = filters.filter((f) => f.field !== "names");
                updateFilters([
                  ...activeFiltersWithoutNames,
                  {
                    field: "names",
                    id: crypto.randomUUID(),
                    operator: id,
                    value: text,
                  },
                ]);
                setOpen(false);
              }}
            />
          ),
        },
        {
          id: "identities",
          label: "Identity",
          shortcut: "i",
          component: (
            <FilterOperatorInput
              label="Identity"
              options={options}
              defaultOption={activeIdentityFilter?.operator}
              defaultText={activeIdentityFilter?.value as string}
              onApply={(id, text) => {
                const activeFiltersWithoutNames = filters.filter((f) => f.field !== "identities");
                updateFilters([
                  ...activeFiltersWithoutNames,
                  {
                    field: "identities",
                    id: crypto.randomUUID(),
                    operator: id,
                    value: text,
                  },
                ]);
                setOpen(false);
              }}
            />
          ),
        },
        {
          id: "keyids",
          label: "Key ID",
          shortcut: "k",
          component: (
            <FilterOperatorInput
              label="Key ID"
              options={keyIdOptions}
              defaultOption={activeKeyIdsFilter?.operator}
              defaultText={activeKeyIdsFilter?.value as string}
              onApply={(id, text) => {
                const activeFiltersWithoutKeyIds = filters.filter((f) => f.field !== "keyIds");
                updateFilters([
                  ...activeFiltersWithoutKeyIds,
                  {
                    field: "keyIds",
                    id: crypto.randomUUID(),
                    operator: id,
                    value: text,
                  },
                ]);
                setOpen(false);
              }}
            />
          ),
        },
        {
          id: "tags",
          label: "Tags",
          shortcut: "t",
          component: (
            <FilterOperatorInput
              label="Tags"
              options={tagsOptions}
              defaultOption={activeTagsFilter?.operator}
              defaultText={activeTagsFilter?.value as string}
              onApply={(id, text) => {
                const activeFiltersWithoutTags = filters.filter((f) => f.field !== "tags");
                updateFilters([
                  ...activeFiltersWithoutTags,
                  {
                    field: "tags",
                    id: crypto.randomUUID(),
                    operator: id,
                    value: text,
                  },
                ]);
                setOpen(false);
              }}
            />
          ),
        },
        {
          id: "outcomes",
          label: "Outcomes",
          shortcut: "o",
          component: <OutcomesFilter onDrawerClose={() => setOpen(false)} />,
        },
      ]}
      activeFilters={filters}
    >
      <div className="group">
        <Button
          variant="ghost"
          className={cn(
            "group-data-[state=open]:bg-gray-4 px-2 rounded-lg",
            filters.length > 0 ? "bg-gray-4" : "",
          )}
          aria-label="Filter logs"
          aria-haspopup="true"
          size="md"
          title="Press 'F' to toggle filters"
        >
          <BarsFilter className="text-accent-9 size-4" />
          <span className="text-accent-12 font-medium text-[13px] max-md:hidden">Filter</span>
          {filters.length > 0 && (
            <div className="bg-gray-7 rounded h-4 px-1 text-[11px] font-medium text-accent-12 text-center flex items-center justify-center">
              {filters.length}
            </div>
          )}
        </Button>
      </div>
    </FiltersPopover>
  );
};
