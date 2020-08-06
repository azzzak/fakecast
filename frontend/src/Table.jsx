import React, { useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

const useStyles = makeStyles({
  head_cell: {
    fontWeight: 'bold'
  }
});

export default function DenseTable(props) {
  const classes = useStyles();

  const [rows, setRows] = React.useState([]);

  useEffect(() => {
    (async () => {
      if (props.podcasts) {
        setRows(props.podcasts);
      }
    })();
  }, [props.podcasts]);

  const handleClick = (event, id) => {
    props.setPodcast(id);
    props.setDrawer(true);
  };

  return (
    <TableContainer>
      <Table size="small" aria-label="podcasts">
        <TableHead>
          <TableRow>
            <TableCell className={classes.head_cell}>#</TableCell>
            <TableCell className={classes.head_cell}>Title</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((row, index) => (
            <TableRow key={row.id} hover
              onClick={(e) => handleClick(e, row.id)}
            >
              <TableCell style={{ width: '20px' }}
              >{rows.length - index}</TableCell>
              <TableCell>{row.title}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer >
  );
}