import React, { useState, useEffect } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  TableSortLabel,
  Paper,
  Box,
  CircularProgress,
  Typography,
  Checkbox,
  IconButton,
  Tooltip,
  useTheme,
} from '@mui/material';
import {
  FilterList as FilterIcon,
  Search as SearchIcon,
} from '@mui/icons-material';

export interface Column<T> {
  id: keyof T | string;
  label: string;
  minWidth?: number;
  align?: 'left' | 'right' | 'center';
  format?: (value: any) => string | React.ReactNode;
  sortable?: boolean;
}

export interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  loading?: boolean;
  error?: string;
  selectable?: boolean;
  onRowClick?: (row: T) => void;
  onSelectionChange?: (selectedRows: T[]) => void;
  defaultSortBy?: keyof T;
  defaultSortDirection?: 'asc' | 'desc';
  rowsPerPageOptions?: number[];
  defaultRowsPerPage?: number;
  showSearch?: boolean;
  showFilters?: boolean;
  emptyStateMessage?: string;
}

function descendingComparator<T>(a: T, b: T, orderBy: keyof T) {
  if (b[orderBy] < a[orderBy]) {
    return -1;
  }
  if (b[orderBy] > a[orderBy]) {
    return 1;
  }
  return 0;
}

function getComparator<T>(
  order: 'asc' | 'desc',
  orderBy: keyof T
): (a: T, b: T) => number {
  return order === 'desc'
    ? (a, b) => descendingComparator(a, b, orderBy)
    : (a, b) => -descendingComparator(a, b, orderBy);
}

function stableSort<T>(array: T[], comparator: (a: T, b: T) => number): T[] {
  const stabilizedThis = array.map((el, index) => [el, index] as [T, number]);
  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0]);
    if (order !== 0) return order;
    return a[1] - b[1];
  });
  return stabilizedThis.map((el) => el[0]);
}

export default function DataTable<T extends { id: string | number }>({
  columns,
  data,
  loading = false,
  error,
  selectable = false,
  onRowClick,
  onSelectionChange,
  defaultSortBy,
  defaultSortDirection = 'asc',
  rowsPerPageOptions = [5, 10, 25],
  defaultRowsPerPage = 10,
  showSearch = false,
  showFilters = false,
  emptyStateMessage = 'No data available',
}: DataTableProps<T>) {
  const theme = useTheme();
  const [order, setOrder] = useState<'asc' | 'desc'>(defaultSortDirection);
  const [orderBy, setOrderBy] = useState<keyof T | string>(defaultSortBy || 'id');
  const [selected, setSelected] = useState<string[]>([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(defaultRowsPerPage);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    if (onSelectionChange) {
      const selectedRows = data.filter((row) =>
        selected.includes(String(row.id))
      );
      onSelectionChange(selectedRows);
    }
  }, [selected, data, onSelectionChange]);

  const handleRequestSort = (property: keyof T | string) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
  };

  const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.checked) {
      const newSelected = data.map((n) => String(n.id));
      setSelected(newSelected);
      return;
    }
    setSelected([]);
  };

  const handleClick = (event: React.MouseEvent<unknown>, id: string) => {
    if (!selectable) return;

    const selectedIndex = selected.indexOf(id);
    let newSelected: string[] = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selected, id);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selected.slice(1));
    } else if (selectedIndex === selected.length - 1) {
      newSelected = newSelected.concat(selected.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selected.slice(0, selectedIndex),
        selected.slice(selectedIndex + 1)
      );
    }

    setSelected(newSelected);
  };

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const isSelected = (id: string) => selected.indexOf(id) !== -1;

  if (error) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          p: 3,
        }}
      >
        <Typography color="error">{error}</Typography>
      </Box>
    );
  }

  const filteredData = searchQuery
    ? data.filter((row) =>
        Object.values(row).some(
          (value) =>
            value &&
            String(value).toLowerCase().includes(searchQuery.toLowerCase())
        )
      )
    : data;

  const sortedData = stableSort(
    filteredData,
    getComparator(order, orderBy as keyof T)
  );

  const paginatedData = sortedData.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  return (
    <Paper sx={{ width: '100%', mb: 2 }}>
      {(showSearch || showFilters) && (
        <Box
          sx={{
            p: 2,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            borderBottom: `1px solid ${theme.palette.divider}`,
          }}
        >
          {showSearch && (
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <SearchIcon sx={{ color: 'action.active', mr: 1 }} />
              <input
                type="text"
                placeholder="Search..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={{
                  border: 'none',
                  outline: 'none',
                  fontSize: '1rem',
                  background: 'transparent',
                }}
              />
            </Box>
          )}
          {showFilters && (
            <Tooltip title="Filter list">
              <IconButton>
                <FilterIcon />
              </IconButton>
            </Tooltip>
          )}
        </Box>
      )}
      <TableContainer>
        <Table stickyHeader aria-label="sticky table">
          <TableHead>
            <TableRow>
              {selectable && (
                <TableCell padding="checkbox">
                  <Checkbox
                    indeterminate={
                      selected.length > 0 && selected.length < data.length
                    }
                    checked={
                      data.length > 0 && selected.length === data.length
                    }
                    onChange={handleSelectAllClick}
                  />
                </TableCell>
              )}
              {columns.map((column) => (
                <TableCell
                  key={String(column.id)}
                  align={column.align}
                  style={{ minWidth: column.minWidth }}
                >
                  {column.sortable !== false ? (
                    <TableSortLabel
                      active={orderBy === column.id}
                      direction={
                        orderBy === column.id ? order : 'asc'
                      }
                      onClick={() => handleRequestSort(column.id)}
                    >
                      {column.label}
                    </TableSortLabel>
                  ) : (
                    column.label
                  )}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell
                  colSpan={
                    columns.length + (selectable ? 1 : 0)
                  }
                  align="center"
                >
                  <CircularProgress />
                </TableCell>
              </TableRow>
            ) : paginatedData.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={
                    columns.length + (selectable ? 1 : 0)
                  }
                  align="center"
                >
                  {emptyStateMessage}
                </TableCell>
              </TableRow>
            ) : (
              paginatedData.map((row) => {
                const isItemSelected = isSelected(String(row.id));
                return (
                  <TableRow
                    hover
                    onClick={(event) => {
                      if (selectable) {
                        handleClick(event, String(row.id));
                      }
                      if (onRowClick) {
                        onRowClick(row);
                      }
                    }}
                    role="checkbox"
                    aria-checked={isItemSelected}
                    tabIndex={-1}
                    key={row.id}
                    selected={isItemSelected}
                    sx={{ cursor: onRowClick ? 'pointer' : 'default' }}
                  >
                    {selectable && (
                      <TableCell padding="checkbox">
                        <Checkbox checked={isItemSelected} />
                      </TableCell>
                    )}
                    {columns.map((column) => {
                      const value = row[column.id as keyof T];
                      return (
                        <TableCell
                          key={String(column.id)}
                          align={column.align}
                        >
                          {column.format
                            ? column.format(value)
                            : value}
                        </TableCell>
                      );
                    })}
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={rowsPerPageOptions}
        component="div"
        count={filteredData.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
      />
    </Paper>
  );
}

// Example usage:
/*
import DataTable, { Column } from '@components/shared/DataTable';

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  status: string;
}

const columns: Column<User>[] = [
  { id: 'name', label: 'Name', minWidth: 170 },
  { id: 'email', label: 'Email', minWidth: 100 },
  {
    id: 'role',
    label: 'Role',
    minWidth: 170,
    align: 'right',
    format: (value) => value.toUpperCase(),
  },
  {
    id: 'status',
    label: 'Status',
    minWidth: 170,
    align: 'right',
    format: (value) => <StatusChip status={value} />,
  },
];

const MyComponent = () => {
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);

  return (
    <DataTable<User>
      columns={columns}
      data={users}
      loading={loading}
      error={error}
      selectable
      onSelectionChange={setSelectedUsers}
      onRowClick={(row) => console.log('clicked row:', row)}
      showSearch
      showFilters
    />
  );
};
*/
