#include <gtk/gtk.h>
#include <stdlib.h>

GtkWidget *list = NULL;
int note = 0;

static void add_note(GtkWidget *widget, gpointer data) {
  g_print("create new note\n");

  char s[20];
  snprintf(s, 20, "%s %d", "new label", ++note);
  GtkWidget *label = gtk_label_new(s);
  gtk_container_add(GTK_CONTAINER(list), label);
  gtk_widget_show_all(list);
}

static void activate(GtkApplication *app, gpointer user_data) {
  GtkWidget *window = gtk_application_window_new(app);
  gtk_window_set_title(GTK_WINDOW(window), "errnote");
  gtk_window_set_default_size(GTK_WINDOW(window), 640, 480);
  gtk_window_set_position(GTK_WINDOW(window), GTK_WIN_POS_CENTER);

  GtkWidget *grid = gtk_grid_new();
  {
    gtk_grid_set_row_spacing(GTK_GRID(grid), 5);
    gtk_grid_set_column_spacing(GTK_GRID(grid), 5);
    gtk_widget_set_margin_left(grid, 5);
    gtk_widget_set_margin_right(grid, 5);
    gtk_widget_set_margin_top(grid, 5);
    gtk_widget_set_margin_bottom(grid, 5);

    gtk_container_add(GTK_CONTAINER(window), grid);
  }

  list = gtk_list_box_new();
  {
    GdkRGBA white;
    white.red = 1;
    white.green = 1;
    white.blue = 1;
    white.alpha = 1;
    gtk_widget_override_background_color(list, GTK_STATE_FLAG_NORMAL, &white);
  }

  GtkWidget *label = gtk_label_new("my label 1");
  GtkWidget *label2 = gtk_label_new("my label 2");
  gtk_container_add(GTK_CONTAINER(list), label);
  gtk_container_add(GTK_CONTAINER(list), label2);
    /* Add the column to the view. */

  GtkWidget *button = gtk_button_new_with_label("New");
  {
    g_signal_connect(button, "clicked", G_CALLBACK(add_note), NULL);
    //g_signal_connect_swapped(button, "clicked", G_CALLBACK(gtk_widget_destroy), window);
  }

  GtkWidget *title = gtk_entry_new();
  {
    gtk_widget_set_hexpand(title, TRUE);
  }

  GtkWidget *view = gtk_text_view_new();
  {
    gtk_widget_set_hexpand(view, TRUE);
    gtk_widget_set_vexpand(view, TRUE);
    //gtk_widget_set_halign(view, gtkAligns[halign]);
    //gtk_widget_set_valign(view, gtkAligns[valign]);

    gtk_container_set_border_width(GTK_CONTAINER(view), 5);
    //gtk_text_view_set_monospace(view, TRUE);

    GtkTextBuffer *buffer = gtk_text_view_get_buffer(GTK_TEXT_VIEW(view));
    gtk_text_buffer_set_text(buffer, "Hello, this is some text", -1);
  }

  gtk_grid_attach(GTK_GRID(grid), button, 0, 0, 1, 1);
  gtk_grid_attach(GTK_GRID(grid), list, 0, 1, 1, 1);
  gtk_grid_attach(GTK_GRID(grid), title, 1, 0, 1, 1);
  gtk_grid_attach(GTK_GRID(grid), view, 1, 1, 1, 1);

  //GtkWidget *button_box = gtk_button_box_new(GTK_ORIENTATION_HORIZONTAL);
  //gtk_container_add(GTK_CONTAINER(window), button_box);

  //gtk_container_add (GTK_CONTAINER(button_box), button);

  gtk_widget_show_all(window);
}

int main(int argc, char **argv) {
  GtkApplication *app = gtk_application_new("org.garbagecollected.errnote", G_APPLICATION_FLAGS_NONE);
  g_signal_connect(app, "activate", G_CALLBACK(activate), NULL);
  int status = g_application_run(G_APPLICATION(app), argc, argv);
  g_object_unref(app);

  return status;
}
