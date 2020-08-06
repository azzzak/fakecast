import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import AttachFileIcon from '@material-ui/icons/AttachFile';

const useStyles = makeStyles(theme => ({
  paper: {
    margin: theme.spacing(1),
    width: theme.spacing(32),
    height: theme.spacing(32),
    backgroundColor: '#dfdfdf',
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center center',
    backgroundSize: 'contain',
  },
  icon: {
    bottom: 3,
    right: 5
  },
  upload: {
    display: 'none',
  },
  label: {
    display: 'none',
  }
}));

export default function Cover(props) {
  const classes = useStyles();

  return (
    <Paper className={classes.paper}
      elevation={2}
      style={{
        backgroundImage: `url(${props.cover})`,
      }}>
      <Grid container justify="flex-end" alignItems="flex-end"
        style={{
          height: '100%',
        }}>

        <input
          accept="image/jpeg, image/png"
          className={classes.upload}
          id="cover"
          name="cover"
          type="file"
          onChange={(e) => {
            props.upload(e);
            e.target.value = null;
          }}
        />

        {props.cover ? <IconButton className={classes.icon}
          onClick={() => props.delete()} >
          <DeleteIcon fontSize="large" />
        </IconButton> : <label htmlFor="cover">
            <IconButton className={classes.icon} component="span">
              <AttachFileIcon fontSize="large" />
            </IconButton>
          </label>}
      </Grid>
    </Paper >
  );
}