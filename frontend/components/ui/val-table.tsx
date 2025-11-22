"use client";

import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TableProps } from '@mui/material';
import { ReactNode } from 'react';

type TableColumn = { id: string; label: ReactNode; align?: 'left' | 'right' | 'center' };
type TableRowData = Record<string, ReactNode>;

type VTableProps = {
  columns: TableColumn[];
  rows: TableRowData[];
} & Omit<TableProps, 'children'>;

export function VTable({ columns, rows, ...rest }: VTableProps) {
  return (
    <TableContainer
      component={Paper}
      sx={{ background: 'var(--valtaris-surface)', border: '1px solid var(--valtaris-border)' }}
    >
      <Table size="small" {...rest}>
        <TableHead>
          <TableRow sx={{ '& th': { color: 'var(--valtaris-text-muted)', borderBottom: '1px solid var(--valtaris-border)' } }}>
            {columns.map((col) => (
              <TableCell key={col.id} align={col.align ?? 'left'}>
                {col.label}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((row, idx) => (
            <TableRow
              key={idx}
              hover
              sx={{
                '& td': { borderBottom: '1px solid var(--valtaris-border)' },
                '&:hover': { backgroundColor: 'var(--valtaris-surface-subtle)' },
              }}
            >
              {columns.map((col) => (
                <TableCell key={col.id} align={col.align ?? 'left'}>
                  {row[col.id]}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
