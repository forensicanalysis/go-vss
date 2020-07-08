#include <stdlib.h>
#include <stdio.h>
#include<fcntl.h> 
#include <unistd.h>
#include <libvshadow.h>
#include <errno.h>

#include <libvshadow_libcerror.h>
#include <libvshadow_volume.h>
#include <libbfio_handle.h>
#include <libbfio_file_io_handle.h>
#include <libcfile_file.h>

#include "_cgo_export.h"

//This wrapper is needed since a Go function can not define the required const keyword
ssize_t read_handle_write_wrapper(
                 intptr_t *io_handle,
                 const uint8_t *buffer,
                 size_t size,
                 libcerror_error_t **error ){
					 uint8_t *buf = buffer;
					return reader_handle_write(io_handle, buf, size, error);
				 }

//Ovewriting internal structures to use specified file descriptor
int libvshadow_volume_open_file_descriptor(
     libvshadow_volume_t *volume,
     int file_descriptor,
     libcerror_error_t **error )

{

	libbfio_handle_t *file_io_handle              = NULL;
	libvshadow_internal_volume_t *internal_volume = NULL;
	static char *function                         = "libvshadow_volume_open";

    int file_io_handle_opened_in_library          = 0;

	if( volume == NULL )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_ARGUMENTS,
		 LIBCERROR_ARGUMENT_ERROR_INVALID_VALUE,
		 "%s: invalid volume.",
		 function );

		return( -1 );
	}
	internal_volume = (libvshadow_internal_volume_t *) volume;

	if( libbfio_file_initialize(
	     &file_io_handle,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_INITIALIZE_FAILED,
		 "%s: unable to create file IO handle.",
		 function );

		goto on_error;
	}

    //_______________________________________________volume_open_file_io_handle_________________________________________________________________
    //   volume, file_io_handle, access_flags, error

    libbfio_internal_handle_t *internal_handle = NULL;
    internal_handle = (libbfio_internal_handle_t *) file_io_handle;

    libbfio_file_io_handle_t *io_handle;
    io_handle = (libbfio_file_io_handle_t *) internal_handle->io_handle;
	
	internal_handle->access_flags = LIBVSHADOW_OPEN_READ;
	internal_handle->open_on_demand = 0;

	io_handle->access_flags = LIBVSHADOW_OPEN_READ;
	io_handle->name= "";
	io_handle->name_size=0;

    libcfile_internal_file_t *internal_file = NULL;
    internal_file = (libcfile_internal_file_t *) io_handle->file;
    internal_file->descriptor = file_descriptor;
	internal_file->access_flags = LIBVSHADOW_OPEN_READ;
	internal_file->current_offset = 0;
	internal_file->size = lseek(file_descriptor, 0, SEEK_END); 

    file_io_handle_opened_in_library = 1;

	if( libvshadow_volume_open_read(
	     internal_volume,
	     file_io_handle,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_IO,
		 LIBCERROR_IO_ERROR_READ_FAILED,
		 "%s: unable to read from file IO handle.",
		 function );

		goto on_error;
	}
#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_grab_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to grab read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif
	internal_volume->file_io_handle                   = file_io_handle;
	internal_volume->file_io_handle_opened_in_library = file_io_handle_opened_in_library;

#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_release_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to release read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif

    //_______________________________________________volume_open_file_io_handle_________________________________________________________________

#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_grab_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to grab read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif
	internal_volume->file_io_handle_created_in_library = 1;

#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_release_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to release read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif
	return( 1 );

on_error:
	if( file_io_handle != NULL )
	{
		libbfio_handle_free(
		 &file_io_handle,
		 NULL );
	}
	return( -1 );
}


//Overwriting internal callback functions to use a go ReaderAt object instead of a file descriptor
int libvshadow_volume_open_go_reader(
     libvshadow_volume_t *volume,
	 int goIndex,
     libcerror_error_t **error){

		 libbfio_handle_t *file_io_handle              = NULL;
	libvshadow_internal_volume_t *internal_volume = NULL;
	static char *function                         = "libvshadow_volume_open";

    int file_io_handle_opened_in_library          = 0;

	if( volume == NULL )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_ARGUMENTS,
		 LIBCERROR_ARGUMENT_ERROR_INVALID_VALUE,
		 "%s: invalid volume.",
		 function );

		return( -1 );
	}
	internal_volume = (libvshadow_internal_volume_t *) volume;

	//This call initializes a file with default functions: they are working on a file descriptor
	if( libbfio_file_initialize(
	     &file_io_handle,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_INITIALIZE_FAILED,
		 "%s: unable to create file IO handle.",
		 function );

		goto on_error;
	}

    //--------------This Section overwrites the default functions of the internal file handle------------------------//
	//-----------------Idea: Call functions, that are defined in Go and operate on a Reader-------------------------//

    libbfio_internal_handle_t *internal_handle = NULL;
    internal_handle = (libbfio_internal_handle_t *) file_io_handle;

#ifdef _WIN32
	internal_handle->io_handle = (long long int*) goIndex;
#else
	internal_handle->io_handle = (long int*) goIndex;
#endif
	internal_handle->free_io_handle = &reader_handle_free;
	internal_handle->clone_io_handle = &reader_handle_clone;
	internal_handle->open = &reader_handle_open;
	internal_handle->close = &reader_handle_close;
	internal_handle->read = &reader_handle_read;
	internal_handle->write = &read_handle_write_wrapper;
	internal_handle->seek_offset = &reader_handle_seak_offset;
	internal_handle->exists = &reader_handle_exists;
	internal_handle->is_open = &reader_handle_is_open;
	internal_handle->get_size = &reader_handle_get_size;


    file_io_handle_opened_in_library = 1;

	
	//This call opens a volume from a file handle
	if( libvshadow_volume_open_read(
	     internal_volume,
	     file_io_handle,
	     error ) != 1 )
	{

		libcerror_error_backtrace_fprint(*error, stdout);
		
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_IO,
		 LIBCERROR_IO_ERROR_READ_FAILED,
		 "%s: unable to read from file IO handle.",
		 function );

		goto on_error;
	}
#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_grab_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to grab read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif
	internal_volume->file_io_handle                   = file_io_handle;
	internal_volume->file_io_handle_opened_in_library = file_io_handle_opened_in_library;

#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_release_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to release read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif

    //_______________________________________________volume_open_file_io_handle_________________________________________________________________

#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_grab_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to grab read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif
	internal_volume->file_io_handle_created_in_library = 1;

#if defined( HAVE_LIBVSHADOW_MULTI_THREAD_SUPPORT )
	if( libcthreads_read_write_lock_release_for_write(
	     internal_volume->read_write_lock,
	     error ) != 1 )
	{
		libcerror_error_set(
		 error,
		 LIBCERROR_ERROR_DOMAIN_RUNTIME,
		 LIBCERROR_RUNTIME_ERROR_SET_FAILED,
		 "%s: unable to release read/write lock for writing.",
		 function );

		return( -1 );
	}
#endif
	return( 1 );

on_error:
	if( file_io_handle != NULL )
	{
		libbfio_handle_free(
		 &file_io_handle,
		 NULL );
	}
	return( -1 );
		 return 2;
}
