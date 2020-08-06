import React from 'react';
import TextField from '@material-ui/core/TextField';

export const field = (c) => {
  return (
    <TextField label={c.label}
      type={c.type}
      size="small"
      variant="outlined"
      value={(c.value).toString()}
      className={c.className}
      onChange={c.onChange}
      error={c.error}
      multiline={!!c.rows}
      rows={c.rows}
      helperText={c.helperText}
    />
  )
}

export const upd = (msg, obj, errors) => {
  return (event, p) => {
    const v = event.target.value;
    var err = {
      ...errors
    };

    if (msg.hasOwnProperty(p))
      v.trim().length === 0 ? err[p] = msg[p] : delete err[p];

    var changed = false;
    if (Object.keys(err).length === 0)
      changed = true

    return {
      updated: {
        ...obj,
        [p]: v
      },
      err: err,
      changed: changed,
    }
  }
}

export const trimmer = (obj) => {
  return Object.assign({}, ...Object.entries(obj).map(([k, v]) => {
    return (typeof v === 'string') ? ({
      [k]: v.trim()
    }) : ({
      [k]: v
    })
  }))
}